package config

import (
	"fmt"
	"net"
	"os"
)

type Config struct {
	PG struct {
		URL      string
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		MaxConn  string `env:"POSTGRES_MAX_CONN"`
	}
	Server struct {
		Port string `env:"SERVER_PORT"`
	}
}

func createPgURL(cfg *Config) string {
	hostPort := net.JoinHostPort(cfg.PG.Host, cfg.PG.Port)
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&pool_max_conns=%s",
		cfg.PG.User,
		cfg.PG.Password,
		hostPort,
		cfg.PG.DB,
		cfg.PG.MaxConn,
	)
}

func New() (*Config, error) {
	cfg := &Config{}

	cfg.PG.Host = os.Getenv("POSTGRES_HOST")
	cfg.PG.Port = os.Getenv("POSTGRES_PORT")
	cfg.PG.DB = os.Getenv("POSTGRES_DB")
	cfg.PG.User = os.Getenv("POSTGRES_USER")
	cfg.PG.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.PG.MaxConn = os.Getenv("POSTGRES_MAX_CONN")
	cfg.PG.URL = createPgURL(cfg)

	cfg.Server.Port = os.Getenv("SERVER_PORT")

	return cfg, nil
}
