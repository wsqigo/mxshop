package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

var (
	client proto.GoodsClient
)

func init() {
	conn, err := grpc.Dial("192.168.2.2:4764", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = proto.NewGoodsClient(conn)
}

func TestGetBrandList(t *testing.T) {
	rsp, err := client.GetBrandList(context.Background(), &proto.BrandFilterRequest{})
	assert.Nil(t, err)
	t.Log(rsp.Total)
	for _, brand := range rsp.Data {
		t.Log(brand.Name)
	}
}
