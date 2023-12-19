package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

func TestGetGoodsList(t *testing.T) {
	rsp, err := client.GetGoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
	})
	assert.Nil(t, err)
	t.Log(rsp.Total)
	for _, category := range rsp.Data {
		t.Log(category.Name)
	}
}
