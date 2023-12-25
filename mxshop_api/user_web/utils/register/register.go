package register

import "mxshop_api/user_web/utils/register/consul"

type Registry interface {
	Register(id, name, address string, port int, tags []string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) Registry {
	return &consul.Registry{
		Host: host,
		Port: port,
	}
}
