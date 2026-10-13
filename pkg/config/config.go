// pkg/config/config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	GeminiAPIKey string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL environment variable is required")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, fmt.Errorf("config: GEMINI_API_KEY environment variable is required")
	}

	return &Config{
		Port:         port,
		DatabaseURL:  dbURL,
		GeminiAPIKey: geminiKey,
	}, nil
}
