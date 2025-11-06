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

	// App
	Port       string
	JWTSecret  string
	AdminEmail string
	AdminPass  string
}

// LoadConfig membaca env vars dan mengembalikannya dalam struct Config
func LoadConfig() *Config {
	// Load file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, membaca dari environment...")
	}

	// Ambil semua env yang sudah di-inject oleh docker-compose
	return &Config{
		DBHost:     getEnv("DB_HOST", "postgres-db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "stokloka"),
		DBPassword: getEnv("DB_PASSWORD", "33f1d019d55113334df954c14a6f6aee4600569c58c3000eb15a836154474bd9f"),
		DBName:     getEnv("DB_NAME", "auth_db"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		Port:       getEnv("PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET_KEY", "xqe1qZ27010tXnXOr7GKmM2Kk4whp+Em/eRwN0HkGU6/APHGtWqgZlNDtOtaDLIIpq01H2vsRZNJ0Jbit1wKw=="),
		AdminEmail: getEnv("DEFAULT_ADMIN_EMAIL", "admin@stokloka.com"),
		AdminPass:  getEnv("DEFAULT_ADMIN_PASS", "admin123"),
	}
}

// getEnv untuk membaca env var dengan nilai default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Environment variable %s tidak ditemukan. Menggunakan fallback: %s", key, fallback)
	return fallback
}
