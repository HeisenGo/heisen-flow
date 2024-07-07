package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func ReadStandard(cfgPath string) (Config, error) {
	var cfg Config

	fullAbsPath, err := absPath(cfgPath)
	if err != nil {
		return cfg, err
	}

	viper.SetConfigFile(fullAbsPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func absPath(cfgPath string) (string, error) {
	if !filepath.IsAbs(cfgPath) {
		return filepath.Abs(cfgPath)
	}
	return cfgPath, nil
}

func MustReadStandard(configPath string) Config {
	cfg, err := ReadStandard(configPath)
	if err != nil {
		panic(err)
	}
	return cfg
}
