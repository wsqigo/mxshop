package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/handler"
	"mxshop_srvs/goods_srv/initialize"
	"mxshop_srvs/goods_srv/proto"
	"mxshop_srvs/goods_srv/utils"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	var err error
	if port == 0 {
		port, err = utils.GetFreePort()
	}
	if err != nil {
		zap.S().Errorw("获取空闲端口失败", "err", err)
		return
	}

	zap.S().Infow("地址端口信息", "port:", port)
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 注册服务健康检查
	// https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	// health.NewServer实现了health接口，不需要单独写个health接口
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic("failed to register service. err: " + err.Error())
	}

	// 生成注册对象
	serviceId := fmt.Sprint(uuid.NewV4())
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceId,
		Port:    port,
		Tags:    global.ServerConfig.Tags,
		Address: global.ServerConfig.Host,
		Check: &api.AgentServiceCheck{ // 对应的检查对象
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "15s",
		},
	}

	// 1. 如何启动两个服务
	// 2. 即使我能通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("register service failed, err: " + err.Error())
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
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Error("注销失败")
	}
	zap.S().Info("注销成功")
}
