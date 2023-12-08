package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo       JWTConfig     `mapstructure:"jwt"`
	AliSmsInfo    AliSmsConfig  `mapstructure:"sms"`
	RedisInfo     RedisConfig   `mapstructure:"redis"`
	ConsulInfo    ConsulConfig  `mapstructure:"consul"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstructure:"expire"`
}

type AliSmsConfig struct {
	ApiKey    string `mapstructure:"key"`
	ApiSecret string `mapstructure:"secret"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
