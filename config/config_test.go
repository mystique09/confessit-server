package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewConfig(t *testing.T) {
  cfgLoader := NewConfigLoader("../", "app")
  serverCfg, err := NewServerConfig(cfgLoader)
  if err != nil {
    t.Error("error creating server config")
  }

  tokenConfig, err := NewTokenConfig(cfgLoader)
  if err != nil {
    t.Error("error creating token config")
  }

  cfg := NewConfig(serverCfg, tokenConfig)
  assert.NotNil(t, cfg)
  assert.Empty(t, cfg.ServerConfig().Host())
}
