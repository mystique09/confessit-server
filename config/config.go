package config

import (
	"cnfs/domain"
	"time"

	"github.com/spf13/viper"
)

type configLoader struct {
	reader *viper.Viper
}

func NewConfigLoader(path, name string) domain.IConfigLoader {
	reader := viper.New()
	reader.AddConfigPath(path)
	reader.SetConfigName(name)
	reader.AutomaticEnv()

	return configLoader{reader}
}

func (loader configLoader) Unmarshal(cfg interface{}) error {
	if err := loader.reader.ReadInConfig(); err != nil {
		return err
	}

	return loader.reader.Unmarshal(&cfg)
}

type config struct {
	server domain.IServerConfig
	token  domain.ITokenConfig
}

func NewConfig(server domain.IServerConfig, token domain.ITokenConfig) domain.IConfig {
	return config{
		server: server,
		token:  token,
	}
}

func (cfg config) ServerConfig() domain.IServerConfig {
	return cfg.server
}

func (cfg config) TokenConfig() domain.ITokenConfig {
	return cfg.token
}

type serverConfig struct {
	Host        string `mapstructure:"HOST"`
	Port        string `mapstructure:"PORT"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	ClientUrl   string `mapstructure:"CLIENT_URL"`
}

func NewServerConfig(loader domain.IConfigLoader) (domain.IServerConfig, error) {
	var cfg serverConfig
	if err := loader.Unmarshal(&cfg); err != nil {
		return serverConfig{}, err
	}
	return cfg, nil
}

func (cfg serverConfig) GetHost() string {
	return cfg.Host
}

func (cfg serverConfig) GetPort() string {
	return cfg.Port
}

func (cfg serverConfig) GetDatabaseUrl() string {
	return cfg.DatabaseUrl
}

func (cfg serverConfig) GetClientUrl() string {
	return cfg.ClientUrl
}

type tokenConfig struct {
	AuthSecretKey         string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	AccessTokenSecretKey  string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenSecretKey string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func NewTokenConfig(loader domain.IConfigLoader) (domain.ITokenConfig, error) {
	var cfg tokenConfig
	if err := loader.Unmarshal(&cfg); err != nil {
		return tokenConfig{}, err
	}
	return cfg, nil
}

func (cfg tokenConfig) GetAuthSecretKey() string {
	return cfg.AuthSecretKey
}

func (cfg tokenConfig) GetAccessTokenSecretKey() string {
	return cfg.AccessTokenSecretKey
}

func (cfg tokenConfig) GetAccessTokenDuration() time.Duration {
	return cfg.AccessTokenDuration
}

func (cfg tokenConfig) GetRefreshTokenSecretKey() string {
	return cfg.RefreshTokenSecretKey
}

func (cfg tokenConfig) GetRefreshTokenDuration() time.Duration {
	return cfg.RefreshTokenDuration
}

// func LoadConfig(path, name string) (Config, error) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName(name)
// 	viper.SetConfigType("env")

// 	viper.AutomaticEnv()

// 	if err := viper.ReadInConfig(); err != nil {
// 		return Config{}, err
// 	}

// 	var config Config
// 	if err := viper.Unmarshal(&config); err != nil {
// 		return Config{}, nil
// 	}

// 	log.Println(config)

// 	return config, nil
// }
