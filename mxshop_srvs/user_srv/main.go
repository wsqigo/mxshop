package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
	"net"
)

func main() {
	var ip string
	var port int
	flag.StringVar(&ip, "ip", "0.0.0.0", "ip地址")
	flag.IntVar(&port, "port", 50051, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("ip:", ip)
	zap.S().Info("port:", port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 注册服务健康检查
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
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      global.ServerConfig.Name,
		Port:    port,
		Tags:    []string{"mxshop", "wsqigo", "user-srv"},
		Address: "192.168.2.2",
		Check: &api.AgentServiceCheck{ // 对应的检查对象
			GRPC:                           "192.168.2.2:50051",
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "15s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("register service failed, err: " + err.Error())
	}

	err = server.Serve(listener)
	if err != nil {
		panic("failed to start grpc: " + err.Error())
	}
}
