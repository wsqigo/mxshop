package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop_api/oss_web/config"
)

var (
	Trans ut.Translator

	ServerConfig config.ServerConfig

	NacosConfig config.NacosConfig
)
