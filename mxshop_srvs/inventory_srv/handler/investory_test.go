package handler

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"mxshop_srvs/inventory_srv/proto"
	"sync"
	"testing"
)

var (
	client proto.InventoryClient
)

func init() {
	conn, err := grpc.Dial("192.168.2.2:8095", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = proto.NewInventoryClient(conn)
}

func TestInventoryServer_SetGoodsInv(t *testing.T) {
	for i := 421; i <= 840; i++ {
		resp, err := client.SetGoodsInv(context.Background(), &proto.GoodsInvInfo{
			GoodsId: int32(i),
			Num:     100,
		})
		assert.Nil(t, err)
		t.Log(resp.String())
	}
	resp, err := client.SetGoodsInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 422,
		Num:     40,
	})

	assert.Nil(t, err)
	t.Log("设置库存成功")
	t.Log(resp.String())
}

func TestInventoryServer_GetGoodsInvDetail(t *testing.T) {
	resp, err := client.GetGoodsInvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
	})

	assert.Nil(t, err)
	t.Log("获取库存成功")
	t.Log(resp.String())
}

func TestInventoryServer_Sell(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 80; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Sell(context.Background(), &proto.SellInfo{
				GoodsInfos: []*proto.GoodsInvInfo{
					{GoodsId: 421, Num: 1},
				},
			})
			if assert.Nil(t, err) {
				fmt.Println("库存扣减成功")
				t.Log(resp.String())
			} else {
				fmt.Println(status.FromError(err))
			}
		}()
	}
	wg.Wait()
}

func TestInventoryServer_Repay(t *testing.T) {
	resp, err := client.Repay(context.Background(), &proto.SellInfo{
		GoodsInfos: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if assert.Nil(t, err) {
		fmt.Println("库存归还成功")
		t.Log(resp.String())
	}
}
