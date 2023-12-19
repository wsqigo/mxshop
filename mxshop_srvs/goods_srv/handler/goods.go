package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct{}

func Model2Response(goods *model.Goods) *proto.GoodsInfoResponse {
	info := &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		Images:          goods.Images,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
	return info
}

func (s *GoodsServer) GetGoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var goodsList []*model.Goods
	localDB := global.DB.Model(model.Goods{})

	if request.KeyWords != "" {
		// 关键词搜索
		localDB = localDB.Where("name like ?", "%"+request.KeyWords+"%")
	}

	if request.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
	}

	if request.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
	}

	if request.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", request.PriceMin)
	}

	if request.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", request.PriceMax)
	}

	if request.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", request.Brand)
	}

	// 通过category查询商品
	if request.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, request.TopCategory); result.Error != nil {
			return nil, status.Errorf(codes.Internal, "获取商品分类失败")
		}

		var subQuery string
		switch category.Level {
		case 1:
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category where parent_category_id = %d)", request.TopCategory)
		case 2:
			subQuery = fmt.Sprintf("select id from category where parent_category_id = %d", request.TopCategory)
		case 3:
			subQuery = fmt.Sprint(request.TopCategory)
		}

		//localDB = localDB.Where("category_id in (?)", subQuery) 不行
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	var count int64
	result := localDB.Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	result = localDB.Preload("Category").Preload("Brand").Scopes(Paginate(request.Pages, request.PagePerNums)).Find(&goodsList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品列表错误")
	}

	res := &proto.GoodsListResponse{Total: count}
	for _, goods := range goodsList {
		res.Data = append(res.Data, Model2Response(goods))
	}

	return res, nil
}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goodsList []*model.Goods
	result := global.DB.Find(&goodsList, info.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "批量获取商品信息失败")
	}

	res := &proto.GoodsListResponse{Total: result.RowsAffected}
	for _, goods := range goodsList {
		res.Data = append(res.Data, Model2Response(goods))
	}

	return res, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	goods := model.Goods{
		CategoryID:      info.CategoryId,
		BrandID:         info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}

	result := global.DB.Create(&goods)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建商品失败")
	}

	return &proto.GoodsInfoResponse{Id: goods.ID}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Goods{}, info.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除商品分类失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	goods := &model.Goods{
		CategoryID:      info.CategoryId,
		BrandID:         info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}

	result := global.DB.Select("*").Updates(goods)
	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	goods := &model.Goods{}

	result := global.DB.First(goods, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品详情失败")
	}

	return Model2Response(goods), nil
}
