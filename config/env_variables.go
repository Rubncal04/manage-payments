package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	PORT                 string
	MONGO_URI            string
	MONGO_DB             string
	TELEGRAM_BOT_TOKEN   string
	TWILIO_ACCOUNT_SID   string
	TWILIO_AUTH_TOKEN    string
	TWILIO_FROM_WHATSAPP string
}

func GetVariables() *EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &EnvVariables{
		PORT:                 os.Getenv("PORT"),
		MONGO_URI:            os.Getenv("MONGO_URI"),
		MONGO_DB:             os.Getenv("MONGO_DB"),
		TELEGRAM_BOT_TOKEN:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		TWILIO_ACCOUNT_SID:   os.Getenv("TWILIO_ACCOUNT_SID"),
		TWILIO_AUTH_TOKEN:    os.Getenv("TWILIO_AUTH_TOKEN"),
		TWILIO_FROM_WHATSAPP: os.Getenv("TWILIO_FROM_WHATSAPP"),
	}
}
