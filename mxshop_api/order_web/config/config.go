package config

type GoodsSrvConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

type ServerConfig struct {
	Name           string         `yaml:"name"`
	Host           string         `yaml:"host"`
	Port           int            `yaml:"port"`
	Tags           []string       `yaml:"tags"`
	GoodsSrvConfig GoodsSrvConfig `yaml:"goods_srv"`
	JWTInfo        JWTConfig      `yaml:"jwt"`
	ConsulInfo     ConsulConfig   `yaml:"consul"`
}

type RedisConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Expire int    `yaml:"expire"`
}

type JWTConfig struct {
	SigningKey string `yaml:"key"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// NacosConfig 读取nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
}
