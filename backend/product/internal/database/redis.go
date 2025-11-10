package database

import (
	"context"
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis membuat koneksi ke Redis
func InitRedis(cfg *config.Config) *redis.Client {
	// Buat alamat Redis
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	// Buat client
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // tidak ada password di setup docker-compose
		DB:       0,  // Gunakan database default
	})

	// Tes koneksi
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Gagal konek ke Redis: %v", err)
	}

	log.Println("Koneksi ke Redis berhasil.")
	return rdb
}
