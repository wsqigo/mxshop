package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mxshop_srvs/userop_srv/proto"
	"testing"
)

var (
	userFavClient proto.UserFavClient
)

func init() {
	conn, err := grpc.Dial("192.168.5.104:2850", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	userFavClient = proto.NewUserFavClient(conn)
}

func TestUserOpServer_GetFavList(t *testing.T) {
	resp, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 5,
	})

	assert.Nil(t, err)
	t.Log(resp.Data)
}
