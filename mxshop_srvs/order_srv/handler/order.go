package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"math/rand"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
	"time"
)

type OrderServer struct{}

func GenerateOrderSn(userId int32) string {
	// 订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(now.UnixNano())
	return fmt.Sprintf("%s%d%d", now.Format("20060102150405"), userId, rand.Intn(90)+10)
}

func (s *OrderServer) GetCartItemList(ctx context.Context, info *proto.UserInfo) (*proto.CartItemListResponse, error) {
	//获取用户的购物车列表
	var carts []*model.ShoppingCart

	result := global.DB.Where(model.ShoppingCart{User: info.Id}).Find(&carts)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取购物车列表失败")
	}

	res := &proto.CartItemListResponse{
		Total: result.RowsAffected,
	}

	for _, cart := range carts {
		res.Data = append(res.Data, &proto.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return res, nil
}

func (s *OrderServer) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车- 合并
	var cart model.ShoppingCart

	result := global.DB.Where(model.ShoppingCart{Goods: request.GoodsId, User: request.UserId}).First(&cart)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Internal, "数据库查询错误")
		}
	}

	if result.RowsAffected == 1 {
		// 如果记录已经存在，则合并购物车记录，更新数量
		cart.Nums += request.Nums
	} else {
		cart.User = request.UserId
		cart.Goods = request.GoodsId
		cart.Nums = request.Nums
		cart.Checked = false
	}

	if err := global.DB.Save(&cart).Error; err != nil {
		return nil, status.Error(codes.Internal, "数据库错误")
	}

	return &proto.ShopCartInfoResponse{Id: cart.ID}, nil
}

func (s *OrderServer) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 更新购物车记录，更新数量和选中状态
	updateField := map[string]any{
		"checked": request.Checked,
	}
	if request.Nums > 0 {
		updateField["nums"] = request.Nums
	}

	result := global.DB.Model(model.ShoppingCart{}).Where(model.ShoppingCart{Goods: request.GoodsId, User: request.UserId}).Updates(updateField)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *OrderServer) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	result := global.DB.Where("goods = ? and user = ?", request.GoodsId, request.UserId).Delete(&model.ShoppingCart{})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除购物车失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *OrderServer) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务（跨微服务）
			3. 库存的扣减 - 访问库存服务（跨微服务）
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/

	var shopCarts []*model.ShoppingCart
	result := global.DB.Where(model.ShoppingCart{User: request.UserId, Checked: true}).Find(&shopCarts)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品")
	}

	goodsNumsMap := make(map[int32]int32, len(shopCarts))
	// 商品 -> 数量
	for _, cart := range shopCarts {
		goodsNumsMap[cart.Goods] = cart.Nums
	}

	// 跨服务调用商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: funk.Keys(goodsNumsMap).([]int32),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "批量获取商品信息失败")
	}

	var orderVal float64
	orderGoodsList := make([]*model.OrderGoods, 0, len(goods.Data))
	goodsInvInfo := make([]*proto.GoodsInvInfo, 0, len(goods.Data))
	for _, data := range goods.Data {
		orderVal += data.ShopPrice * float64(goodsNumsMap[data.Id])
		// 还没创建订单，订单id还不知道，等创建好赋值
		orderGoodsList = append(orderGoodsList, &model.OrderGoods{
			Goods:      data.Id,
			GoodsName:  data.Name,
			GoodsImage: data.GoodsFrontImage,
			GoodsPrice: data.ShopPrice,
			Nums:       goodsNumsMap[data.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: data.Id,
			Num:     goodsNumsMap[data.Id],
		})
	}

	// 跨服务调用库存微服务进行库存扣减
	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfos: goodsInvInfo,
	})
	if err != nil {
		return nil, err
	}

	// 生成订单表
	// 20210308xxxx
	tx := global.DB.Begin()
	order := model.OrderInfo{
		User:         request.UserId,
		OrderSn:      GenerateOrderSn(request.UserId),
		OrderMount:   orderVal,
		Address:      request.Address,
		SignerName:   request.Name,
		SignerMobile: request.Mobile,
		Post:         request.Post,
	}

	result = tx.Create(&order)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	for _, orderGoods := range orderGoodsList {
		orderGoods.Order = order.ID
	}

	// 批量插入orderGoodsList
	result = tx.CreateInBatches(orderGoodsList, 100)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	// 删除购物车记录
	result = tx.Unscoped().Where(model.ShoppingCart{User: request.UserId, Checked: true}).Delete(&model.ShoppingCart{})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	tx.Commit()

	return &proto.OrderInfoResponse{Id: order.ID, OrderSn: order.OrderSn, Total: order.OrderMount}, nil
}

func (s *OrderServer) GetOrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 使用场景 1.后台管理系统 2.电商系统查询
	// request.UserId 为默认值，并不会添加判断
	var total int64
	result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{User: request.UserId}).Count(&total)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "数据库内部错误")
	}

	var orders []model.OrderInfo
	result = global.DB.Scopes(Paginate(request.Pages, request.PagePerNums)).
		Where(model.OrderInfo{User: request.UserId}).Find(&orders)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "数据库内部错误")
	}

	res := &proto.OrderListResponse{Total: total}
	orderList := make([]*proto.OrderInfoResponse, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, Model2Response(order))
	}

	res.Data = orderList
	return res, nil
}

func (s *OrderServer) GetOrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo

	// 这个订单的id是否是当前用户的订单，如果在web层用户传递过来一个id的订单
	// web 层应该先查询一下订单id是否是当前用户的
	// 在个人中心可以这样做，但是如果是后台管理系统，那么只传递order的id
	// 如果是电商系统还需要一个用户的id
	qh := global.DB
	if request.UserId != 0 {
		qh = qh.Where("user = ?", request.UserId)
	}

	result := qh.First(&order, request.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "数据库错误")
	}

	res := &proto.OrderInfoDetailResponse{
		OrderInfo: Model2Response(order),
	}

	// 订单的商品
	var orderGoods []model.OrderGoods
	result = global.DB.Where(model.OrderGoods{Order: order.ID}).Find(&orderGoods)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	for _, goods := range orderGoods {
		res.GoodsItems = append(res.GoodsItems, &proto.OrderItemResponse{
			Id:         goods.Goods,
			OrderId:    goods.Order,
			GoodsId:    goods.Goods,
			GoodsName:  goods.GoodsName,
			GoodsImage: goods.GoodsImage,
			GoodsPrice: goods.GoodsPrice,
			Nums:       goods.Nums,
		})
	}

	return res, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	// 先查询，再更新 实际上有两条sql执行
	// select 和 update 语句
	result := global.DB.Model(model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).
		Update("status", req.Status)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}
