package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası bulunamadı: ", err)
	}

	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}
}
