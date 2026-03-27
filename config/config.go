package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	DSN      string `json:"dsn"`
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: .env файл не найден!")
	}
	godotenv.Load()
	return Config{
		Email:    os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
		Address:  os.Getenv("ADDRESS"),
		DSN:      os.Getenv("DSN"),
	}
}
