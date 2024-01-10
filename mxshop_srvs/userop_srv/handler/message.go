package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/model"
	"mxshop_srvs/userop_srv/proto"
)

func (s *UserOpServer) GetMessageList(ctx context.Context, request *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var messages []*model.LeavingMessage

	result := global.DB.Where(model.LeavingMessage{User: request.UserId}).Find(&messages)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "获取信息列表失败")
	}

	res := &proto.MessageListResponse{
		Total: result.RowsAffected,
	}

	for _, message := range messages {
		res.Data = append(res.Data, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	return res, nil
}

func (s *UserOpServer) CreateMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	message := &model.LeavingMessage{
		User:        request.UserId,
		MessageType: request.MessageType,
		Subject:     request.Subject,
		Message:     request.Message,
		File:        request.File,
	}

	result := global.DB.Create(message)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}

	return &proto.MessageResponse{Id: message.ID}, nil
}
