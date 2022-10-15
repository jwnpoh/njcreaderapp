package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port string
	DSN  string
}

func LoadConfig() (Config, error) {
	godotenv.Load(".env")

	cfg := Config{
		Port: os.Getenv("PORT"),
		DSN:  os.Getenv("PSCALE"),
	}

	return cfg, nil
}
