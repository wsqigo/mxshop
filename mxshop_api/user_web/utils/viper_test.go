package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// 如何将线上和线下的配置文件隔离

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name      string      `mapstructure:"name"`
	MysqlInfo MysqlConfig `mapstructure:"mysql"`
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	// 刚才设置的环境变量想要生效 我们必须得重启goland
}

func TestViperDemo(t *testing.T) {
	isDebug := GetEnvInfo("MXSHOP_DEBUG")
	fmt.Println(isDebug)
	configFileNamePrefix := "config"
	var configFileName string
	if isDebug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("%s-pro.yaml", configFileName)
	}

	cfg := ServerConfig{}
	v := viper.New()
	v.SetConfigFile(configFileName)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	// viper的功能 -- 动态监听变化
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			_ = v.ReadInConfig() //读取配置数据
			_ = v.Unmarshal(&cfg)

		})
	}()
	time.Sleep(time.Second * 3000)
}
