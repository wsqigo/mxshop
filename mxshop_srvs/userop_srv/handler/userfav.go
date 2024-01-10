package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/model"
	"mxshop_srvs/userop_srv/proto"
)

func (s *UserOpServer) GetFavList(ctx context.Context, request *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var userFavList []*model.UserFav

	// 查询用户的收藏记录
	// 查询某件商品被哪些用户收藏了
	result := global.DB.Where(model.UserFav{User: request.UserId, Goods: request.GoodsId}).Find(&userFavList)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "获取收藏列表失败")
	}

	res := &proto.UserFavListResponse{
		Total: result.RowsAffected,
	}

	for _, userFav := range userFavList {
		res.Data = append(res.Data, &proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}

	return res, nil
}

func (s *UserOpServer) AddUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	userFav := &model.UserFav{
		User:  request.UserId,
		Goods: request.GoodsId,
	}

	result := global.DB.Create(&userFav)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建收藏失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *UserOpServer) DeleteUserFav(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	// todo: 测试 userId 为空，是否有where判断
	result := global.DB.Unscoped().Where("goods=? and user=?", request.GoodsId, request.UserId).Delete(&model.UserFav{})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除收藏失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *UserOpServer) GetUserFavDetail(ctx context.Context, request *proto.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav
	result := global.DB.Where("goods=? and user=?", request.GoodsId, request.UserId).Find(&userFav)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
		}
		return nil, status.Errorf(codes.Internal, "查询收藏失败")
	}
	return &emptypb.Empty{}, nil
}
