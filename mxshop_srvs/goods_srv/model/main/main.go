package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/goods_srv/model"
	"os"
	"time"
)

func main() {
	dsn := "root:root@tcp(192.168.136.130)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("建立数据库连接失败")
	}

	db.AutoMigrate(&model.Category{}, &model.Brand{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})
}
