package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl           string        `mapstructure:"DATABASE_URL"`
	Host                  string        `mapstructure:"HOST"`
	PasetoSymmetricKey    string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	AccessTokenSecretKey  string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenSecretKey string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ClientUrl             string        `mapstructure:"CLIENT_URL"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath("../")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, nil
	}

	return config, nil
}
