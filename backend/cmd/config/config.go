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

type SheetsConfig struct {
	Credentials string
	SheetID     string
	SheetName   string
}

type Config struct {
	Port string
	DSN  string
	MailServiceConfig
	SheetsConfig
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
		SheetsConfig: SheetsConfig{
			Credentials: os.Getenv("SHEET_CREDENTIALS"),
			SheetID:     os.Getenv("SHEET_ID"),
			SheetName:   os.Getenv("SHEET_NAME"),
		},
	}

	return cfg, nil
}
