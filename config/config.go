package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl           string        `mapstructure:"DATABASE_URL"`
	Host                  string        `mapstructure:"HOST"`
	Port                  string        `mapstructure:"PORT"`
	PasetoSymmetricKey    string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	AccessTokenSecretKey  string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenSecretKey string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ClientUrl             string        `mapstructure:"CLIENT_URL"`
}

func LoadConfig(path, name string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, nil
	}

	log.Println(config)

	return config, nil
}
