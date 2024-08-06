package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", nil
	}
	return value, nil
}
