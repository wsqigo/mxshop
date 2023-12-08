package config

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul"`
}
