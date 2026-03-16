package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	PORT                      string
	AllowedOrigins            []string
	MONGO_URL                 string
	JWT_SECRET                string
	JWT_TOKEN_EXPIRATION_TIME time.Duration
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatal("MONGO_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	allowedOrigins := parseOrigins(os.Getenv("ALLOWED_ORIGINS"))
	if len(allowedOrigins) == 0 {
		log.Fatal("ALLOWED_ORIGINS environment variable is required")
	}

	jwtExpirationStr := os.Getenv("JWT_TOKEN_EXPIRATION_TIME")
	if jwtExpirationStr == "" {
		log.Fatal("JWT_TOKEN_EXPIRATION_TIME environment variable is required")
	}

	jwtExpiration, err := time.ParseDuration(jwtExpirationStr)
	if err != nil {
		log.Fatalf("Invalid JWT_TOKEN_EXPIRATION_TIME: %v", err)
	}

	return &Config{
		PORT:                      port,
		AllowedOrigins:            allowedOrigins,
		MONGO_URL:                 mongoURL,
		JWT_SECRET:                jwtSecret,
		JWT_TOKEN_EXPIRATION_TIME: jwtExpiration,
	}
}

func parseOrigins(origins string) []string {

	parts := strings.Split(origins, ",")

	var cleaned []string
	for _, o := range parts {
		o = strings.TrimSpace(o)
		if o != "" {
			cleaned = append(cleaned, o)
		}
	}

	return cleaned
}
