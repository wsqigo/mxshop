package handler

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mxshop_srvs/user_srv/proto"
	"testing"
)

var userClient proto.UserClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic("grpc connect failed, err: " + err.Error())
	}

	userClient = proto.NewUserClient(conn)
}

func TestUserServer_GetUserList(t *testing.T) {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		PNum:  1,
		PSize: 2,
	})

	assert.Nil(t, err)
	for _, user := range rsp.Data {
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})

		assert.Nil(t, err)
		assert.True(t, checkRsp.Success)
	}
}

func TestUserServer_CreateUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("wsqigo%d", i),
			Mobile:   fmt.Sprintf("1912415529%d", i),
			Password: "admin123",
		})
		assert.Nil(t, err)
		fmt.Println(rsp.Id)
	}
}
