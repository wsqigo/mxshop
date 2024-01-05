package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct{}

func (s *InventoryServer) SetGoodsInv(ctx context.Context, info *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	inventory := model.Inventory{}
	result := global.DB.Where(model.Inventory{Goods: info.GoodsId}).First(&inventory)

	inventory.Goods = info.GoodsId
	inventory.Stocks = info.Num

	result = global.DB.Save(&inventory)
	if result.Error != nil {
		zap.S().Errorf("set goods inventory failed: %v", result.Error)
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) GetGoodsInvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	inventory := model.Inventory{}
	result := global.DB.Where(&model.Inventory{Goods: info.GoodsId}).First(&inventory)
	if result.Error != nil {
		zap.S().Errorf("get goods inventory detail failed: %v", result.Error)
		return nil, result.Error
	}

	return &proto.GoodsInvInfo{
		GoodsId: inventory.Goods,
		Num:     inventory.Stocks,
	}, nil
}

func (s *InventoryServer) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10, 2:5, 3:20]
	// 数据库基本的一个应用场景：数据库事务
	tx := global.DB.Begin()
	for _, goodsInfo := range info.GoodsInfos {
		inv := model.Inventory{}

		fmt.Println("开始获取锁")
		mutex := global.RedSync.NewMutex(fmt.Sprintf("goods_%d", goodsInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "lock failed")
		}
		result := global.DB.Where(model.Inventory{Goods: goodsInfo.GoodsId}).First(&inv)
		if result.Error != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}

		// 判断库存是否充足
		if inv.Stocks < goodsInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减，会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodsInfo.Num

		// inv 有 primary key，不需要 where 判断
		result = tx.Model(&inv).Update("stocks", inv.Stocks)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "unlock failed")
		}
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) Repay(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 库存归还: 1. 订单超时归还 2. 订单创建失败，归还之前扣减的库存 3. 手动归还
	tx := global.DB.Begin()

	for _, goodsInfo := range info.GoodsInfos {
		inv := model.Inventory{}
		result := global.DB.Where(model.Inventory{Goods: goodsInfo.GoodsId}).First(&inv)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}

		// 扣减，会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks += goodsInfo.Num

		result = tx.Model(&inv).Update("stocks", inv.Stocks)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}
