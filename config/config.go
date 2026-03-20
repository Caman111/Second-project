package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func LoadConfig() Config {
	godotenv.Load()
	return Config{
		Email:    os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
		Address:  os.Getenv("ADDRESS"),
	}
}
