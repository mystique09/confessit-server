package domain

import "time"

type IConfig interface {
  ServerConfig() IServerConfig
  AuthSecretKey() string
  TokenConfig() ITokenConfig
}

type ITokenConfig interface {
  AccessTokenSecrerKey() string
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
