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
  token domain.ITokenConfig
}

func NewConfig(server domain.IServerConfig, token domain.ITokenConfig) domain.IConfig {
  return config{
    server: server,
    token: token,
  }
}

func (cfg config) ServerConfig() domain.IServerConfig {
  return cfg.server
}

func (cfg config) TokenConfig() domain.ITokenConfig {
  return cfg.token
}

type serverConfig struct {
  host string `mapstructure:"HOST"`
  port string `mapstructure:"PORT"`
  databaseUrl string `mapstructure:"DATABASE_URL"`
  clientUrl string `mapstructure:"CLIENT_URL"`
}

func NewServerConfig(loader domain.IConfigLoader) (domain.IServerConfig, error) {
  var cfg serverConfig
  if err := loader.Unmarshal(&cfg); err != nil {
    return serverConfig{}, err
  }
  return cfg, nil
}

func (cfg serverConfig) Host() string {
  return cfg.host
}

func (cfg serverConfig) Port() string {
  return cfg.port
}

func (cfg serverConfig) DatabaseUrl() string {
  return cfg.databaseUrl
}

func (cfg serverConfig) ClientUrl() string {
  return cfg.clientUrl
}

type tokenConfig struct {
  authSecretKey string `mapstructure:"PASETO_SYMMETRIC_KEY"`
  accessTokenSecretKey string `mapstructure:"ACCESS_TOKEN_SECRET"`
  accessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
  refreshTokenSecretKey string `mapstructure:"REFRESH_TOKEN_SECRET"`
  refreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func NewTokenConfig(loader domain.IConfigLoader) (domain.ITokenConfig, error) {
  var cfg tokenConfig
  if err := loader.Unmarshal(&cfg); err != nil {
    return tokenConfig{}, err
  }
  return cfg, nil
}

func (cfg tokenConfig) AuthSecretKey() string {
  return cfg.authSecretKey
}

func (cfg tokenConfig) AccessTokenSecretKey() string {
  return cfg.accessTokenSecretKey
}

func(cfg tokenConfig) AccessTokenDuration() time.Duration {
  return cfg.accessTokenDuration
}

func (cfg tokenConfig) RefreshTokenSecretKey() string {
  return cfg.refreshTokenSecretKey
}

func (cfg tokenConfig) RefreshTokenDuration() time.Duration {
  return cfg.refreshTokenDuration
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
