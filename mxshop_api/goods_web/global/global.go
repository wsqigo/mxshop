package global

import (
	"mxshop_api/goods_web/config"
	"mxshop_api/goods_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator

	ServerConfig config.ServerConfig

	NacosConfig config.NacosConfig

	GoodsSrvClient proto.GoodsClient
)
