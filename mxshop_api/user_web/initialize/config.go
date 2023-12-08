package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_api/user_web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	// 设置的环境变量想要生效，必须得重启goland
	return viper.GetBool(env)
}

func InitConfig() {
	isDebug := GetEnvInfo("MXSHOP_DEBUG")

	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_web/%s-pro.yaml", configFilePrefix)
	if isDebug {
		configFileName = fmt.Sprintf("user_web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic("init config failed, err: " + err.Error())
	}

	// 配置信息应该为全局变量
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic("unable to decode into struct, err: " + err.Error())
	}
	zap.S().Infof("配置信息: %v", global.ServerConfig)

	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			zap.S().Info("config file changed:", e.Name)
			if err := v.Unmarshal(&global.ServerConfig); err != nil {
				panic("unable to decode into struct, err: " + err.Error())
			}
			zap.S().Infof("配置信息: %v", global.ServerConfig)
		})
	}()
}
