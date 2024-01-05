package handler

import (
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
)

func Paginate[T constraints.Integer](pageNum, pageSize T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize

		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

func Model2Response(order model.OrderInfo) *proto.OrderInfoResponse {
	// 在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	// 这里要搞清，哪些字段是有默认值
	orderInfoRsp := &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SignerMobile,
		AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return orderInfoRsp
}
