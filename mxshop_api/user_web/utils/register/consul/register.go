package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	return client.Agent().ServiceDeregister(serviceId)
}

func (r *Registry) Register(id, name, address string, port int, tags []string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic("failed to register service. err: " + err.Error())
	}

	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check: &api.AgentServiceCheck{ // 对应的检查对象
			HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "15s",
		},
	}

	return client.Agent().ServiceRegister(registration)
}
