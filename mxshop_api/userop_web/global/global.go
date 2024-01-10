package global

import (
	"mxshop_api/userop_web/config"
	"mxshop_api/userop_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator

	ServerConfig config.ServerConfig

	NacosConfig config.NacosConfig

	GoodsSrvClient proto.GoodsClient

	MessageSrvClient proto.MessageClient
	AddressSrvClient proto.AddressClient
	UserFavSrvClient proto.UserFavClient
)
