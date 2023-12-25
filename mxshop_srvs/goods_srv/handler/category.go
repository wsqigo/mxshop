package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

// GetAllCategoryList 商品分类
func (s *GoodsServer) GetAllCategoryList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	// 可以放在web层处理
	var categoryList []*model.Category
	result := global.DB.Where("level = ?", 1).Preload("SubCategoryList.SubCategoryList").Find(&categoryList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品分类列表失败")
	}

	bytes, err := json.Marshal(categoryList)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal错误")
	}

	result = global.DB.Find(&categoryList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取商品分类列表失败")
	}

	res := &proto.CategoryListResponse{}
	categoryInfoList := make([]*proto.CategoryInfoResponse, 0, result.RowsAffected)
	for _, category := range categoryList {
		categoryInfoList = append(categoryInfoList, &proto.CategoryInfoResponse{
			Id:             category.ID,
			Name:           category.Name,
			ParentCategory: category.ParentCategoryID,
			Level:          category.Level,
			IsTab:          category.IsTab,
		})
	}

	res.Total = result.RowsAffected
	res.Data = categoryInfoList
	res.JsonData = string(bytes)
	return res, nil
}

func (s *GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	category := model.Category{}
	if result := global.DB.First(&category, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	info := &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategoryList []*model.Category
	result := global.DB.Where(&model.Category{ParentCategoryID: request.Id}).Find(&subCategoryList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取子分类错误")
	}

	subCategoryInfoList := make([]*proto.CategoryInfoResponse, 0, result.RowsAffected)
	for _, subCategory := range subCategoryList {
		subCategoryInfoList = append(subCategoryInfoList, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})
	}

	return &proto.SubCategoryListResponse{
		Total:           1,
		Info:            info,
		SubCategoryList: subCategoryInfoList,
	}, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := &model.Category{
		Name:             request.Name,
		Level:            request.Level,
		IsTab:            request.IsTab,
		ParentCategoryID: request.ParentCategory,
	}

	result := global.DB.Create(category)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建商品分类失败")
	}

	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Category{}, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除商品分类失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := &model.Category{
		Name:             request.Name,
		ParentCategoryID: request.ParentCategory,
		Level:            request.Level,
		IsTab:            request.IsTab,
	}

	result := global.DB.Select("*").Where("id = ?", request.Id).Updates(category)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新商品分类失败")
	}
	return &emptypb.Empty{}, nil
}
