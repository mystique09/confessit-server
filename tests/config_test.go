package config

import (
	"cnfs/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewConfig(t *testing.T) {
	cfgLoader := config.NewConfigLoader("../", "app")
	serverCfg, err := config.NewServerConfig(cfgLoader)
	if err != nil {
		t.Error("error creating server config")
	}

	tokenConfig, err := config.NewTokenConfig(cfgLoader)
	if err != nil {
		t.Error("error creating token config")
	}

	cfg := config.NewConfig(serverCfg, tokenConfig)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.ServerConfig().GetHost())
}
