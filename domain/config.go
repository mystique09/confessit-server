package domain

import "time"

type IConfigLoader interface {
  Unmarshal(interface{}) error
}

type IConfig interface {
	ServerConfig() IServerConfig
	TokenConfig() ITokenConfig
}

type ITokenConfig interface {
	AuthSecretKey() string
	AccessTokenSecretKey() string
	AccessTokenDuration() time.Duration
	RefreshTokenSecretKey() string
	RefreshTokenDuration() time.Duration
}

type IServerConfig interface {
	Host() string
	Port() string
	DatabaseUrl() string
	ClientUrl() string
}
