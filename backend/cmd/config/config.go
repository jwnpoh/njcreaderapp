package config

import (
	"github.com/joho/godotenv"
	"os"
)

type MailServiceConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Config struct {
	Port string
	DSN  string
	MailServiceConfig
}

func LoadConfig() (Config, error) {
	godotenv.Load(".env")

	cfg := Config{
		Port: os.Getenv("PORT"),
		DSN:  os.Getenv("PSCALE"),
		MailServiceConfig: MailServiceConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASS"),
		},
	}

	return cfg, nil
}
