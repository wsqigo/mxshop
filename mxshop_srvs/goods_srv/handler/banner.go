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

// GetBannerList 轮播图
func (s *GoodsServer) GetBannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	res := &proto.BannerListResponse{}

	var banners []model.Banner
	result := global.DB.Find(&banners)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取轮播图列表失败")
	}

	bannerInfoList := make([]*proto.BannerResponse, 0, result.RowsAffected)
	for _, banner := range banners {
		bannerInfoList = append(bannerInfoList, &proto.BannerResponse{
			Id:    banner.ID,
			Index: banner.Index,
			Image: banner.Image,
			Url:   banner.Url,
		})
	}

	res.Total = result.RowsAffected
	res.Data = bannerInfoList
	return res, nil
}

func (s *GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerInfoRequest) (*proto.BannerResponse, error) {
	banner := &model.Banner{
		Image: request.Image,
		Index: request.Index,
		Url:   request.Url,
	}

	result := global.DB.Create(banner)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建轮播图失败")
	}

	return &proto.BannerResponse{Id: banner.ID}, nil
}

func (s *GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerInfoRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Banner{}, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "删除轮播图失败")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerInfoRequest) (*emptypb.Empty, error) {
	var banner model.Banner
	if result := global.DB.First(&banner, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	banner.Url = request.Url
	banner.Image = request.Image
	banner.Index = request.Index

	result := global.DB.Where("id = ?", request.Id).Updates(banner)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败")
	}
	return &emptypb.Empty{}, nil
}
