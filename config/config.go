package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	USER_SERVICE    string
	PRODUCT_SERVICE string
	DB_HOST         string
	DB_PORT         string
	DB_USER         string
	DB_PASSWORD     string
	DB_NAME         string
	SIGNING_KEY     string
}

func Load() Config {
	if err := godotenv.Load("C:/imtixon/Product Service/.env"); err != nil {
		log.Print("No .env file found?")
	}

	config := Config{}
	config.DB_HOST = cast.ToString(Coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(Coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(Coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(Coalesce("DB_PASSWORD", "salom"))
	config.DB_NAME = cast.ToString(Coalesce("DB_NAME", "products"))
	config.PRODUCT_SERVICE = cast.ToString(Coalesce("PRODUCT_SERVICE", "50053"))
	config.USER_SERVICE = cast.ToString(Coalesce("USER_SERVICE", "localhost:50051"))
	config.SIGNING_KEY = cast.ToString(Coalesce("SIGNING_KEY", "secret"))

	return config
}

func Coalesce(key string, defaultValue interface{}) interface{} {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
