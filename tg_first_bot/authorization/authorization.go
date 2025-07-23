package authorization

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func MustToken() (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("can't read .env file")
	}

	t := os.Getenv("tg_first_bot")
	if t == "" {
		return "", errors.New("token is empty")
	}
	return t, nil
}
