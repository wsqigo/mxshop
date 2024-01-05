package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mxshop_srvs/order_srv/proto"
	"testing"
)

var (
	client proto.OrderClient
)

func init() {
	conn, err := grpc.Dial("192.168.2.2:5668", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = proto.NewOrderClient(conn)
}

func TestOrderServer_CreateCartItem(t *testing.T) {
	_, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		Nums:    1,
		GoodsId: 421,
	})

	assert.Nil(t, err)
}

func TestOrderServer_GetCartItemList(t *testing.T) {
	resp, err := client.GetCartItemList(context.Background(), &proto.UserInfo{
		Id: 1,
	})
	assert.Nil(t, err)
	for _, item := range resp.Data {
		t.Log(item.Id, item.GoodsId, item.Nums)
	}
}

func TestOrderServer_UpdateCartItem(t *testing.T) {
	resp, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 421,
		Checked: true,
	})
	assert.Nil(t, err)
	t.Log(resp)
}

func TestOrderServer_CreateOrder(t *testing.T) {
	_, err := client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "长沙市",
		Name:    "wsqigo",
		Mobile:  "19124155294",
		Post:    "请尽快发货",
	})
	assert.Nil(t, err)
}

func TestOrderServer_GetOrderDetail(t *testing.T) {
	resp, err := client.GetOrderDetail(context.Background(), &proto.OrderRequest{
		Id: 1,
	})
	assert.Nil(t, err)
	t.Log(resp)
	for _, item := range resp.GoodsItems {
		t.Log(item.GoodsName)
	}
}

func TestOrderServer_GetOrderList(t *testing.T) {
	resp, err := client.GetOrderList(context.Background(), &proto.OrderFilterRequest{})
	assert.Nil(t, err)
	for _, order := range resp.Data {
		t.Log(order.OrderSn)
	}
}
