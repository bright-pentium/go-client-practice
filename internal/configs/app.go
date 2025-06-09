package configs

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port             int
	LogLevel         string
	DbURL            string
	MaxConn          int
	MinConn          int
	SecretKey        string
	SecretExpiration int
	Issuer           string
}

func LoadConfig(envFilePath string) (*AppConfig, error) {
	envMap, err := godotenv.Read(envFilePath)
	if err != nil {
		log.Fatalf("Error reading .env file from %s: %v", envFilePath, err)
	}

	rawPort := getEnv(envMap, "PORT", "8000")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, err
	}

	rawMaxConn := getEnv(envMap, "MAX_CONN", "15")
	maxConn, err := strconv.Atoi(rawMaxConn)
	if err != nil {
		return nil, err
	}

	rawMinConn := getEnv(envMap, "MIN_CONN", "5")
	minConn, err := strconv.Atoi(rawMinConn)
	if err != nil {
		return nil, err
	}

	secretKey := getEnv(envMap, "SECRET_KEY", "")
	if secretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY cannot be empty.")
	}

	issuer := getEnv(envMap, "ISSUER", "")
	if issuer == "" {
		return nil, fmt.Errorf("ISSUER cannot be empty.")
	}

	raawExpiration := getEnv(envMap, "SECRET_EXPIRATION", "900")
	expiration, err := strconv.Atoi(raawExpiration)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Port:             port,
		MaxConn:          maxConn,
		MinConn:          minConn,
		LogLevel:         getEnv(envMap, "LOG_LEVEL", "INFO"),
		SecretKey:        secretKey,
		Issuer:           issuer,
		SecretExpiration: expiration,
	}, nil
}
