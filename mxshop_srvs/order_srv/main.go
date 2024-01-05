package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
	"mxshop_srvs/order_srv/initialize"
	"mxshop_srvs/order_srv/proto"
	"mxshop_srvs/order_srv/utils"
	"mxshop_srvs/order_srv/utils/registry"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvClients()

	port, err := utils.GetFreePort()
	if err != nil {
		zap.S().Panicw("获取空闲端口失败", "err", err)
	}

	zap.S().Infow("地址端口信息", "port:", port)
	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 注册服务健康检查
	// https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	// health.NewServer实现了health接口，不需要单独写个health接口
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	registerClient := registry.NewRegistryClient(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
	)

	serviceId := uuid.NewV4().String()
	err = registerClient.Register(
		serviceId,
		global.ServerConfig.Name,
		global.ServerConfig.Host,
		port,
		global.ServerConfig.Tags,
	)
	if err != nil {
		zap.S().Panicf("服务注册失败：%v", err)
	}

	go func() {
		err = server.Serve(listener)
		if err != nil {
			panic("failed to start grpc: " + err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Error("注销失败")
	}
	zap.S().Info("注销成功")
}
