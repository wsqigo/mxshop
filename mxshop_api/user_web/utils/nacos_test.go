package utils

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func TestNacos(t *testing.T) {
	serverConfigs := []constant.ServerConfig{
		{
			Scheme: "http",
			IpAddr: "192.168.136.130",
			Port:   8848,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "4084e7fa-9c12-4b15-a6d2-a498eb8e979c",
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
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web",
		Group:  "STAGE",
	})

	if err != nil {
		panic(err)
	}

	fmt.Print(content)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web",
		Group:  "STAGE",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件变化")
			fmt.Println("group:"+group+", dataId:", dataId+", data:"+data)
		},
	})
	time.Sleep(3000 * time.Second)
}
