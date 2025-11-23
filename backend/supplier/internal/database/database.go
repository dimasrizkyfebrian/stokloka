package database

import (
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase membuat koneksi, migrasi, dan seeder
func InitDatabase(cfg *config.Config) *gorm.DB {
	// Buat DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	// Buka Koneksi GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Gagal konek ke database %s: %v", cfg.DBName, err)
	}
	log.Printf("Koneksi ke database %s berhasil.", cfg.DBName)

	// Jalankan Migrasi
	runMigrations(db)

	// Jalankan Seeder
	SeedData(db)

	return db
}

// runMigrations menjalankan GORM AutoMigrate
func runMigrations(db *gorm.DB) {
	log.Println("Menjalankan migrasi database...")
	err := db.AutoMigrate(
		&model.Supplier{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}
	log.Println("Migrasi database berhasil.")
}

// SeedData mengisi data master awal
func SeedData(db *gorm.DB) {
	log.Println("Mengecek data seeder...")

	var count int64
	// Cek apakah tabel supplier sudah ada isinya
	db.Model(&model.Supplier{}).Count(&count)

	if count == 0 {
		log.Println("Menjalankan seeder untuk Supplier...")

		suppliers := []model.Supplier{
			{
				KodePemasok:  "SUP-001",
				NamaPemasok:  "PT Pemasok Sejahtera",
				KontakPerson: "Budi Santoso",
				Email:        "budi@pemasoksejahtera.com",
				Telepon:      "081234567890",
				Alamat:       "Jl. Industri Raya No. 123, Jakarta",
			},
		}

		// Hook BeforeCreate di model akan otomatis mengisi ID (UUID)
		if err := db.Create(&suppliers).Error; err != nil {
			log.Printf("Gagal seeding supplier: %v", err)
		} else {
			log.Println("Seeder supplier berhasil dijalankan.")
		}
	} else {
		log.Println("Data supplier sudah ada, skip seeding.")
	}
}
