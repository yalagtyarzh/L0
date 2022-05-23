package main

import (
	"log"

	"github.com/yalagtyarzh/L0/internal/application"
	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/pkg/logging"
)

func main() {
	log.Println("reading environment")
	cfg := config.GetConfig()

	log.Println("logging initializing")
	logger := logging.InitLogger(cfg.InProduction, cfg.AppConfig.LogLevel)
	defer logger.Sync()

	a, err := application.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("running application")
	a.Run()
}
