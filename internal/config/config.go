package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Service  Service
	Platform Platform
	Postgres Postgres
}

type Service struct {
	Port string `env:"AUTH_PORT"`
}

type Platform struct {
	Token string `env:"TOKEN"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

func NewConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg
}
