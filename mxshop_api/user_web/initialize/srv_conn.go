package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/proto"
)

func InitSrvConn() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic("new consul client failed. err: " + err.Error())
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvConfig.Name))
	if err != nil {
		panic("get service user-srv failed, err: " + err.Error())
	}

	userSrvHost := data["user-srv"].Address
	userSrvPort := data["user-srv"].Port
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	// 拨号连接用户grpc服务 跨域的问题 -- 后端解决 也可以前端来解决
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
		return
	}

	// 生成grpc的client并调用接口
	// 可能存在的问题
	// 1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做
	// 4. 一个连接多个goroutine公用，性能 - 连接池 grpc-go-pool
	global.UserSrvClient = proto.NewUserClient(conn)
}
