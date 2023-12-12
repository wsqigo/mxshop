package config

type UserSrvConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

type ServerConfig struct {
	Name          string        `yaml:"name"`
	Port          int           `yaml:"port"`
	UserSrvConfig UserSrvConfig `yaml:"user_srv"`
	JWTInfo       JWTConfig     `yaml:"jwt"`
	AliSmsInfo    AliSmsConfig  `yaml:"sms"`
	RedisInfo     RedisConfig   `yaml:"redis"`
	ConsulInfo    ConsulConfig  `yaml:"consul"`
}

type RedisConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Expire int    `yaml:"expire"`
}

type AliSmsConfig struct {
	ApiKey    string `yaml:"key"`
	ApiSecret string `yaml:"secret"`
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
