package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/proto"
)

func InitSrvClients() {
	// 初始化第三方微服务的client
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=wsqigo", consulInfo.Host, consulInfo.Port,
			global.ServerConfig.GoodsSrvInfo.Name),
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

	// 初始化库存服务连接
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=wsqigo", consulInfo.Host, consulInfo.Port,
			global.ServerConfig.InventorySrvInfo.Name),
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
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}
