package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_srvs/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// InitConfig 从配置文件中读取出对应的配置
func InitConfig() {
	isDebug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_srv/%s-pro.yaml", configFilePrefix)
	if isDebug {
		configFileName = fmt.Sprintf("user_srv/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic("init config failed, err: " + err.Error())
	}

	err = v.Unmarshal(&global.ServerConfig)
	if err != nil {
		panic("unable to decode into struct, err: " + err.Error())
	}

	zap.S().Infof("配置信息: %v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("config file changed:", e.Name)
		err := v.ReadInConfig()
		if err != nil {
			panic("init config failed, err: " + err.Error())
		}

		err = v.Unmarshal(&global.ServerConfig)
		if err != nil {
			panic("unable to decode into struct, err: " + err.Error())
		}
	})
}
