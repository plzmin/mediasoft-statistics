package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"mediasoft-statistics/internal/app"
	"mediasoft-statistics/internal/config"
	"mediasoft-statistics/pkg/logger"
)

func main() {
	log := logger.New()

	cfg := config.Config{}
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatal("failed to retrieve env variables %v", err)
	}

	if err := app.Run(log, cfg); err != nil {
		log.Fatal("error running gateway server %v", err)
	}
}
