package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"mxshop_srvs/inventory_srv/config"
)

var (
	DB *gorm.DB

	RedSync *redsync.Redsync

	ServerConfig config.ServerConfig

	NacosConfig config.NacosConfig
)
