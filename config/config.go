package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	SessionSecret  string
	JWTSecret      string
	JWTExpiresHour int
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	jwtExpires, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "fiber_coa"),
		SessionSecret:  getEnv("SESSION_SECRET", "supersecretkey"),
		JWTSecret:      getEnv("JWT_SECRET", "supersecretkey"),
		JWTExpiresHour: jwtExpires,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
