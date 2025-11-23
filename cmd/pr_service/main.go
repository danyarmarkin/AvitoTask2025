package main

import (
	"AvitoTask2025/config"
	"AvitoTask2025/internal/app"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		log.Fatalf("can not get application config: %s", err)
	}

	logger, err := zap.NewDevelopment()

	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("can not initialize logger: %s", err)
	}

	app.Run(logger, cfg)
}
