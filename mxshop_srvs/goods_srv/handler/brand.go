package handler

import (
	"context"
	"fmt"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

func (g GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	res := &proto.BrandListResponse{}

	var brands []model.Brand
	result := global.DB.Find(&brands)
	fmt.Println(result.RowsAffected)

	brandResList := make([]*proto.BrandInfoResponse, 0, result.RowsAffected)
	for _, brand := range brandResList {
		brandResList = append(brandResList, &proto.BrandInfoResponse{
			Id:   brand.Id,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	res.Total = result.RowsAffected
	res.Data = brandResList
	return res, nil
}
