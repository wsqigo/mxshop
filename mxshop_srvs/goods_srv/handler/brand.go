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

// GetBrandList 品牌和轮播图
func (g *GoodsServer) GetBrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	res := &proto.BrandListResponse{}

	var brands []model.Brand
	result := global.DB.Scopes(Paginate(request.Pages, request.PagePerNums)).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询品牌总数
	var count int64
	global.DB.Find(&model.Brand{}).Count(&count)

	brandInfoList := make([]*proto.BrandInfoResponse, 0, result.RowsAffected)
	for _, brand := range brands {
		brandInfoList = append(brandInfoList, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	res.Total = count
	res.Data = brandInfoList
	return res, nil
}

// CreateBrand 新建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandInfoRequest) (*proto.BrandInfoResponse, error) {
	if result := global.DB.First(&model.Brand{Name: request.Name}); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "品牌已存在")
	}

	brand := &model.Brand{Name: request.Name, Logo: request.Logo}

	global.DB.Save(brand)

	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}

func (s *GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandInfoRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Brand{}, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除品牌失败")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandInfoRequest) (*emptypb.Empty, error) {
	if result := global.DB.First(&model.Brand{Name: request.Name}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	brand := &model.Brand{
		Name: request.Name,
		Logo: request.Logo,
	}

	result := global.DB.Updates(&brand)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败")
	}
	return &emptypb.Empty{}, nil
}
