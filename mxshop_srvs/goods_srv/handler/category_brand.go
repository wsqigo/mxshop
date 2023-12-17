package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

func (s *GoodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var categoryBrands []*model.GoodsCategoryBrand

	var count int64
	result := global.DB.Model(&model.GoodsCategoryBrand{}).Count(&count)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品分类-品牌列表失败")
	}

	result = global.DB.Preload("Category").Preload("Brand").Scopes(Paginate(request.Pages, request.PagePerNums)).Find(&categoryBrands)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品分类-品牌列表失败")
	}

	categoryBrandList := make([]*proto.CategoryBrandResponse, 0, result.RowsAffected)
	for _, categoryBrand := range categoryBrands {
		categoryBrandList = append(categoryBrandList, &proto.CategoryBrandResponse{
			Id: categoryBrand.ID,
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.BrandID,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.CategoryID,
				Name:           categoryBrand.Category.Name,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
			},
		})
	}

	res := &proto.CategoryBrandListResponse{
		Total: count,
		Data:  categoryBrandList,
	}

	return res, nil
}

// GetCategoryBrandList 根据商品分类获取Brand信息
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var categoryBrands []*model.GoodsCategoryBrand
	result := global.DB.Preload("Brand").Where("category_id = ?", request.Id).Find(&categoryBrands)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取Brand信息列表失败")
	}

	brandInfoList := make([]*proto.BrandInfoResponse, 0)
	for _, categoryBrand := range categoryBrands {
		brandInfoList = append(brandInfoList, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brand.ID,
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		})
	}

	res := &proto.BrandListResponse{
		Total: result.RowsAffected,
		Data:  brandInfoList,
	}

	return res, nil
}

func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	categoryBrand := &model.GoodsCategoryBrand{
		CategoryID: request.CategoryId,
		BrandID:    request.BrandId,
	}

	result := global.DB.Create(categoryBrand)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建商品分类-品牌失败")
	}

	return &proto.CategoryBrandResponse{Id: categoryBrand.ID}, nil
}

func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.GoodsCategoryBrand{}, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除商品分类-品牌失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	categoryBrand := &model.GoodsCategoryBrand{
		CategoryID: request.CategoryId,
		BrandID:    request.BrandId,
	}

	result := global.DB.Where("id = ?", request.Id).Updates(categoryBrand)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新商品分类失败")
	}
	return &emptypb.Empty{}, nil
}
