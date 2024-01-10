package handler

import (
	"context"
	"fmt"
	"testing"

	"mxshop_srvs/userop_srv/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	addrClient proto.AddressClient
)

func init() {
	conn, err := grpc.Dial("192.168.5.104:2850", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	addrClient = proto.NewAddressClient(conn)
}

func TestUserOpServer_GetAddressList(t *testing.T) {
	resp, err := addrClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 5,
	})

	assert.Nil(t, err)
	fmt.Println(resp.Data)
}
