package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user-srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}
