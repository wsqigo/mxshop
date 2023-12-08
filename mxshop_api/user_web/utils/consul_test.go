package utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
)

func Register(id string, address string, port int, name string, tags []string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.136.130:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			HTTP:                           "http://192.168.2.2:8021/health",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func TestConsul(t *testing.T) {
	//Register("user-web", "192.168.2.2", 8021, "user-web", []string{"mxshop", "wsqigo"})
	AllServices()
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.136.130:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}

	for key, _ := range data {
		fmt.Println(key)
	}
}

func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.136.130:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}
