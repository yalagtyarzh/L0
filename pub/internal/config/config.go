package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the config for parsing env variables into application
type Config struct {
	STAN struct {
		Cluster string `env:"WB_NATS_CLUSTER" env-default:"test-cluster"`
		Client  string `env:"WB_NATS_CLIENT" env-default:"publisher"`
		Delay   int    `env:"WB_NATS_DELAY" env-default:"5"`
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
				fmt.Println("environment is not OK")
				os.Exit(1)
			}
		},
	)

	return cfg
}
