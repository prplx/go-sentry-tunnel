package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DSN             []string      `env:"DSN, env-required"`
	AllowOrigins    []string      `env:"ALLOW_ORIGINS" env-default:"*"`
	Env             string        `env:"ENV" env-default:"production"`
	Port            string        `env:"PORT" env-default:"3001"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	RequestTimeout  time.Duration `env:"REQUEST_TIMEOUT" env-default:"5s"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
