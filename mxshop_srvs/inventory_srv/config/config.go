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
	Name        string       `yaml:"name"`
	Host        string       `yaml:"host"`
	Tags        []string     `yaml:"tags"`
	MysqlInfo   MysqlConfig  `yaml:"mysql"`
	ConsulInfo  ConsulConfig `yaml:"consul"`
	RedisConfig RedisConfig  `yaml:"redis"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
}
