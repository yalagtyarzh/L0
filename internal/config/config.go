package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	PostgreSQL struct {
		Username string `env:"PSQL_USERNAME" env-required:"true"`
		Password string `env:"PSQL_PASSWORD" env-required:"true"`
		Host     string `env:"PSQL_HOST" env-required:"true"`
		Port     string `env:"PSQL_PORT" env-required:"true"`
		Database string `env:"PSQL_DATABASE" env-required:"true"`
	}
}

func GetConfig() *Config {
	log.Println("reading environment")

	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		fmt.Println("cfg is not OK")
		os.Exit(1)
	}

	return cfg
}
