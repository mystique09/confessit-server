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
	GetAuthSecretKey() string
	GetAccessTokenSecretKey() string
	GetAccessTokenDuration() time.Duration
	GetRefreshTokenSecretKey() string
	GetRefreshTokenDuration() time.Duration
}

type IServerConfig interface {
	GetHost() string
	GetPort() string
	GetDatabaseUrl() string
	GetClientUrl() string
}
