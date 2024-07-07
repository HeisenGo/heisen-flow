package config

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
}

type ServerConfig struct {
	HTTPPort           int    `mapstructure:"http_port"`
	Host               string `mapstructure:"host"`
	TokenExpMinutes    int    `mapstructure:"token_exp_minutes"`
	RefreshTokenExpMin int    `mapstructure:"refresh_token_exp_minute"`
	TokenSecret        string `mapstructure:"token_secret"`
}

type DBConfig struct {
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	DBName string `mapstructure:"db_name"`
}
