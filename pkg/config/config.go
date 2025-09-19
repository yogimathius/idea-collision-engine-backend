package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	OpenAIAPIKey     string
	StripeSecretKey  string
	Environment      string
	CORSOrigins      []string
	RateLimitRPS     int
	CacheExpiration  int // seconds
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	rateLimitRPS, _ := strconv.Atoi(getEnvWithDefault("RATE_LIMIT_RPS", "10"))
	cacheExpiration, _ := strconv.Atoi(getEnvWithDefault("CACHE_EXPIRATION", "300"))

	config := &Config{
		Port:             getEnvWithDefault("PORT", "8080"),
		DatabaseURL:      getEnvWithDefault("DATABASE_URL", ""),
		RedisURL:         getEnvWithDefault("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:        getEnvWithDefault("JWT_SECRET", "your-secret-key-change-in-production"),
		OpenAIAPIKey:     getEnvWithDefault("OPENAI_API_KEY", ""),
		StripeSecretKey:  getEnvWithDefault("STRIPE_SECRET_KEY", ""),
		Environment:      getEnvWithDefault("ENVIRONMENT", "development"),
		CORSOrigins:      []string{getEnvWithDefault("CORS_ORIGINS", "http://localhost:5173")},
		RateLimitRPS:     rateLimitRPS,
		CacheExpiration:  cacheExpiration,
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}
	if c.StripeSecretKey == "" && c.Environment == "production" {
		return fmt.Errorf("STRIPE_SECRET_KEY is required in production")
	}
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}