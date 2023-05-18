package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	App struct {
		Port        int
		ServiceName string
	}

	Log struct {
		Level string
	}

	Redis struct {
		Url  string
		Pass string
	}

	Session struct {
		TTL time.Duration
	}

	UsersService struct {
		URI string
	}
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
