package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct{}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Model2Response(user model.User) *proto.UserInfoResponse {
	// 在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	// 这里要搞清，哪些字段是有默认值
	userInfoRsp := &proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     int64(user.Role),
	}

	if user.Birthday != nil {
		userInfoRsp.Birthday = user.Birthday.Unix()
	}

	return userInfoRsp
}

func (s *UserServer) GetUserList(ctx context.Context, info *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = result.RowsAffected

	global.DB.Scopes(Paginate(int(info.PNum), int(info.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := Model2Response(user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}

	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 通过手机号码查询用户
	var user model.User

	res := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IDRequest) (*proto.UserInfoResponse, error) {
	// 通过 ID 查询用户
	var user model.User

	res := global.DB.First(&user, req.Id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, info *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 新建用户
	var user model.User

	res := global.DB.Where(&model.User{Mobile: info.Mobile}).First(&user)
	if res.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = info.Mobile
	user.NickName = info.NickName

	// 密码加密
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encodedPwd := password.Encode(info.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, info *proto.UpdateUserInfo) (*empty.Empty, error) {
	// 个人中心更新用户
	var user model.User
	res := global.DB.First(&user, info.Id)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthday := time.Unix(info.Birthday, 0)
	user.NickName = info.NickName
	user.Birthday = &birthday
	user.Gender = info.Gender
	res = global.DB.Updates(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, info *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	// 校验密码
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	passwordInfo := strings.Split(info.EncryptedPassword, "$")
	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
