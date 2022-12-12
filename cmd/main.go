package main

import (
	"cnfs/config"
	"cnfs/handler"

	"github.com/labstack/gommon/log"
)

func main() {
	cfg, err := config.LoadConfig(".", "app")
	if err != nil {
		log.Fatal(err)
	}

	handler.Launch(&cfg)
}
