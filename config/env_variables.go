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
	JWT_SECRET_KEY       string
	REDIS_ADDRESS        string
	REDIS_PASSWORD       string
	REDIS_PORT           string
	REDIS_DATABASES      string
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
		JWT_SECRET_KEY:       os.Getenv("JWT_SECRET_KEY"),
		REDIS_ADDRESS:        os.Getenv("REDIS_ADDRESS"),
		REDIS_PASSWORD:       os.Getenv("REDIS_PASSWORD"),
		REDIS_PORT:           os.Getenv("REDIS_PORT"),
		REDIS_DATABASES:      os.Getenv("REDIS_DATABASES"),
	}
}
