package config

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ServerConfig struct {
	Name       string       `yaml:"name"`
	Host       string       `yaml:"host"`
	Tags       []string     `yaml:"tags"`
	MysqlInfo  MysqlConfig  `yaml:"mysql"`
	RedisInfo  RedisConfig  `yaml:"redis"`
	ConsulInfo ConsulConfig `yaml:"consul"`

	// 商品微服务的配置
	GoodsSrvInfo     GoodsSrvConfig   `yaml:"goods_srv"`
	InventorySrvInfo InventorySrvInfo `yaml:"inventory_srv"`
}

type GoodsSrvConfig struct {
	Name string `yaml:"name"`
}

type InventorySrvInfo struct {
	Name string `yaml:"name"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
}
