package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	DiscordAppID string
	DatabaseURL  string
	Environment  string
	LogLevel     string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		DiscordAppID: os.Getenv("DISCORD_APP_ID"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Environment:  getEnvOrDefault("APP_ENV", "dev"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}
	if cfg.DiscordAppID == "" {
		return nil, fmt.Errorf("DISCORD_APP_ID is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
