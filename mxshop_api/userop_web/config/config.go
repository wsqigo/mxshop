package config

type SrvConfig struct {
	Name string `yaml:"name"`
}

type ServerConfig struct {
	Name string   `yaml:"name"`
	Host string   `yaml:"host"`
	Port int      `yaml:"port"`
	Tags []string `yaml:"tags"`

	GoodsSrvConfig  SrvConfig    `yaml:"goods_srv"`
	UserOpSrvConfig SrvConfig    `yaml:"userop_srv"`
	JWTInfo         JWTConfig    `yaml:"jwt"`
	ConsulInfo      ConsulConfig `yaml:"consul"`
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

type AlipayConfig struct {
	AppId        string `yaml:"app_id"`
	PrivateKey   string `yaml:"private_key"`
	AliPublicKey string `yaml:"ali_public_key"`
	NotifyUrl    string `yaml:"notify_url"`
	ReturnUrl    string `yaml:"return_url"`
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
