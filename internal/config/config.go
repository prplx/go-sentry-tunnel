package config

import (
	"fmt"
	"os"
	"time"

	"github.com/prplx/go-sentry-tunnel/internal/errors"
)

type Config struct {
	Env             string
	Port            string
	ShutdownTimeout time.Duration
}

func MustLoad() *Config {
	config := &Config{
		Env:             os.Getenv("ENV"),
		Port:            os.Getenv("PORT"),
		ShutdownTimeout: 10 * time.Second,
	}

	if config.Env == "" {
		panic(fmt.Errorf("%w: ENV", errors.ErrorEnvVariableRequired))
	}

	if config.Port == "" {
		panic(fmt.Errorf("%w: PORT", errors.ErrorEnvVariableRequired))
	}

	return config
}
