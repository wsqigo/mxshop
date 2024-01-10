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

func (s *UserOpServer) GetAddressList(ctx context.Context, request *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []*model.Address

	result := global.DB.Where(model.Address{User: request.UserId}).Find(&addresses)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "获取地址列表失败")
	}

	res := &proto.AddressListResponse{
		Total: result.RowsAffected,
	}

	for _, address := range addresses {
		res.Data = append(res.Data, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}

	return res, nil
}

func (s *UserOpServer) CreateAddress(ctx context.Context, request *proto.AddressRequest) (*proto.AddressResponse, error) {
	address := &model.Address{
		User:         request.UserId,
		Province:     request.Province,
		City:         request.City,
		District:     request.District,
		Address:      request.Address,
		SignerName:   request.SignerName,
		SignerMobile: request.SignerMobile,
	}

	result := global.DB.Create(address)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}

	return &proto.AddressResponse{Id: address.ID}, nil
}

func (s *UserOpServer) DeleteAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	// todo: 测试 userId 为空，是否有where判断
	result := global.DB.Where("id=? and user=?", request.Id, request.UserId).Delete(&model.Address{})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *UserOpServer) UpdateAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	var address model.Address
	result := global.DB.Where("id=? and user=?", request.Id, request.UserId).First(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
		}
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}

	updateField := map[string]any{}
	if address.Province != "" {
		address.Province = request.Province
	}

	if address.City != "" {
		address.City = request.City
	}

	if address.District != "" {
		address.District = request.District
	}

	if address.Address != "" {
		address.Address = request.Address
	}

	if address.SignerName != "" {
		address.SignerName = request.SignerName
	}

	if address.SignerMobile != "" {
		address.SignerMobile = request.SignerMobile
	}

	if err := global.DB.Model(&model.Address{}).Updates(updateField).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "更新失败: %v", err)
	}

	return &emptypb.Empty{}, nil
}
