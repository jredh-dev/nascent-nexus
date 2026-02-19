package config

import (
	"os"
)

// Config holds all configuration for the calendar service.
type Config struct {
	Port   string
	DBPath string
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:   envOr("CAL_PORT", "8085"),
		DBPath: envOr("CAL_DB_PATH", "cal.db"),
	}
}
