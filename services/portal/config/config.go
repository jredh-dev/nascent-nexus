package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Giveaway GiveawayConfig
	Session  SessionConfig
}

// GiveawayConfig holds giveaway service settings.
type GiveawayConfig struct {
	DBPath string // path to giveaway SQLite database file
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string
	Env  string
}

// DBConfig holds database settings.
type DBConfig struct {
	Path string // path to SQLite database file
}

// SessionConfig holds session/cookie settings.
type SessionConfig struct {
	Secret string // HMAC key for signing session cookies
	MaxAge int    // session duration in seconds (default: 7 days)
}

// Load returns application configuration from environment variables.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "portal.db"),
		},
		Giveaway: GiveawayConfig{
			DBPath: getEnv("GIVEAWAY_DB_PATH", "giveaway.db"),
		},
		Session: SessionConfig{
			Secret: getEnv("SESSION_SECRET", ""),
			MaxAge: getEnvInt("SESSION_MAX_AGE", 604800), // 7 days
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
