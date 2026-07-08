package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	Port              string
	OxeliaGatewayMode bool
	EncryptionKey     string
}

var Cfg *Config

func Load() *Config {
	_ = godotenv.Load()

	Cfg = &Config{
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		Port:              getEnv("MUSICBOX_PORT", "8003"),
		OxeliaGatewayMode: getEnvBool("OXELIA_GATEWAY_MODE", false),
		EncryptionKey:     getEnv("ENCRYPTION_KEY", "default-encryption-key-change-me-in-prod-32b"),
	}

	if Cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if len(Cfg.EncryptionKey) < 32 {
		log.Fatal("ENCRYPTION_KEY must be at least 32 bytes")
	}

	return Cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
