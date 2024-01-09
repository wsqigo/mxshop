package initialize

import (
	"fmt"

	"mxshop_api/order_web/global"
	"mxshop_api/order_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=wsqigo", consulInfo.Host, consulInfo.Port,
			global.ServerConfig.GoodsSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Panicf("[InitSrvConn] 连接商品服务失败")
	}

	// 生成grpc的client并调用接口
	// 可能存在的问题
	// 1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做
	// 4. 一个连接多个goroutine公用，性能 - 连接池 grpc-go-pool
	global.GoodsSrvClient = proto.NewGoodsClient(conn)

	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=wsqigo", consulInfo.Host, consulInfo.Port,
			global.ServerConfig.OrderSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Panicf("[InitSrvConn] 连接订单服务失败")
	}
	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=wsqigo", consulInfo.Host, consulInfo.Port,
			global.ServerConfig.InventoryConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Panicf("[InitSrvConn] 连接库存服务失败")
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}
