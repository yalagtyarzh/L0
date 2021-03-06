package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// STAN struct holds information about STAN connection
type STAN struct {
	Cluster string `env:"WB_NATS_CLUSTER" env-default:"test-cluster"`
	Client  string `env:"WB_NATS_CLIENT" env-default:"XXD-Client"`
	Channel string `env:"WB_NATS_CHANNEL" env-default:"XD-CHANNEL"`
	Delay   int    `env:"WB_NATS_DELAY" env-default:"5"`
}

// Config is the config for parsing env variables into application
type Config struct {
	UseCache     bool `env:"WB_USE_CACHE" env-default:"false"`
	InProduction bool `env:"WB_IN_PROD" env-default:"false"`
	AppConfig    struct {
		LogLevel string `env:"WB_LOG_LEVEL" env-default:"debug"`
	}
	STAN
	Listen struct {
		IP   string `env:"WB_IP" env-default:"127.0.0.1"`
		Port string `env:"WB_PORT" env-default:":8080"`
	}
	PostgreSQL struct {
		Username string `env:"WB_PSQL_USERNAME" env-required:"true"`
		Password string `env:"WB_PSQL_PASSWORD" env-required:"true"`
		Host     string `env:"WB_PSQL_HOST" env-required:"true"`
		Port     string `env:"WB_PSQL_PORT" env-required:"true"`
		Database string `env:"WB_PSQL_DATABASE" env-required:"true"`
	}
}

var cfg *Config
var once sync.Once

// GetConfig gets environment variables from system
func GetConfig() *Config {
	once.Do(
		func() {
			cfg = &Config{}

			if err := cleanenv.ReadEnv(cfg); err != nil {
				fmt.Printf("environment is not OK: %s\n", err)
				os.Exit(1)
			}
		},
	)

	return cfg
}
