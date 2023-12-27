package config

type JWTConfig struct {
	SigningKey string `yaml:"key"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type OssConfig struct {
	ApiKey      string `yaml:"api_key"`
	ApiSecret   string `yaml:"api_secret"`
	Host        string `yaml:"host"`
	CallBackUrl string `yaml:"call_back_url"`
	UploadDir   string `yaml:"upload_dir"`
}

type ServerConfig struct {
	Name       string       `yaml:"name"`
	Host       string       `yaml:"host"`
	Port       int          `yaml:"port"`
	Tags       []string     `yaml:"tags"`
	JWTInfo    JWTConfig    `yaml:"jwt"`
	ConsulInfo ConsulConfig `yaml:"consul"`
	OssInfo    OssConfig    `yaml:"oss"`
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
