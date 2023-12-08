package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop_api/user_web/config"
	"mxshop_api/user_web/proto"
)

var (
	ServerConfig  config.ServerConfig
	Trans         ut.Translator
	UserSrvClient proto.UserClient
)
