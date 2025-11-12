package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi aplikasi
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost string
	RedisPort string

	// RabbitMQ
	RabbitHost string
	RabbitPort string

	// App
	Port      string
	JWTSecret string
}

// LoadConfig membaca env vars dan mengembalikannya dalam struct Config
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, membaca dari environment...")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "stokloka"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "product_db"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		RabbitHost: getEnv("RABBITMQ_HOST", "localhost"),
		RabbitPort: getEnv("RABBITMQ_PORT", "5672"),

		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET_KEY", "default_secret"),
	}
}

// getEnv adalah helper untuk membaca env var dengan nilai default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Environment variable %s tidak ditemukan. Menggunakan fallback: %s", key, fallback)
	return fallback
}
