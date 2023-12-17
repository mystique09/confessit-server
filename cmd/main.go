package main

import (
	"cnfs/config"
	"cnfs/handler"

	"log"
)

func main() {
	cfgLoader := config.NewConfigLoader(".", "app")
	serveConfig, err := config.NewServerConfig(cfgLoader)
	if err != nil {
		log.Fatal(err)
	}

	tokenConfig, err := config.NewTokenConfig(cfgLoader)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig(serveConfig, tokenConfig)

	handler.Launch(cfg)
}
