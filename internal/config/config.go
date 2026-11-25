// Package config lê a configuração a partir de variáveis de ambiente.
package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("variável de ambiente DATABASE_URL não definida")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{DatabaseURL: url, Port: port}, nil
}
