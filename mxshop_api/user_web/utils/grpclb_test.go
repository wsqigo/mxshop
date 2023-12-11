package utils

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mxshop_api/user_web/proto"
	"testing"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func TestLoadBalance(t *testing.T) {
	conn, err := grpc.Dial(
		"consul://192.168.136.130:8500/user-srv?wait=14s&tag=wsqigo",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		client := proto.NewUserClient(conn)
		rsp, err := client.GetUserList(context.Background(), &proto.PageInfo{
			PNum:  1,
			PSize: 2,
		})
		if err != nil {
			panic(err)
		}
		for index, data := range rsp.Data {
			fmt.Print(index, data)
		}
	}
}
