package config

import (
	"github.com/joho/godotenv"
	"os"
)

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

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
	TelegramConfig
}

func LoadConfig() (Config, error) {
	godotenv.Load(".env")

	cfg := Config{
		Port: os.Getenv("PORT"),
		DSN:  os.Getenv("COCKROACH"),
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
		TelegramConfig: TelegramConfig{
			BotToken: os.Getenv("BOT_TOKEN"),
			ChatID:   os.Getenv("CHAT_ID"),
		},
	}

	return cfg, nil
}
