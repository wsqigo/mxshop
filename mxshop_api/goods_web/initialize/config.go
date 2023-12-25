package initialize

import (
	"fmt"

	"mxshop_api/goods_web/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	// 设置的环境变量想要生效，必须得重启goland
	return viper.GetBool(env)
}

func InitConfig() {
	isDebug := GetEnvInfo("MXSHOP_DEBUG")

	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods_web/%s-prod.yaml", configFilePrefix)
	if isDebug {
		configFileName = fmt.Sprintf("goods_web/%s-stage.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic("init config failed, err: " + err.Error())
	}

	// 配置信息应该为全局变量
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic("unable to decode into struct, err: " + err.Error())
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	// 从nacos中读取配置信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
			Scheme: "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: serverConfigs,
		ClientConfig:  &clientConfig,
	})
	if err != nil {
		zap.S().Panic("init nacos client fail. err:", err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Panic("read nacos config fail, err", err)
	}

	err = yaml.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Panicf("failed to unmarshal yaml: %v", err)
	}
}
