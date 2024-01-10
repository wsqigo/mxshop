package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mxshop_srvs/userop_srv/proto"
	"testing"
)

var (
	messageClient proto.MessageClient
)

func init() {
	conn, err := grpc.Dial("192.168.2.2:2850", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	messageClient = proto.NewMessageClient(conn)
}

func TestUserOpServer_GetMessageList(t *testing.T) {
	resp, err := messageClient.GetMessageList(context.Background(), &proto.MessageRequest{
		UserId: 5,
	})

	assert.Nil(t, err)
	t.Log(resp.Data)
}
