package database

import (
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
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
		cfg.DBName, // Ini akan berisi 'product_db'
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
		&model.Category{},
		&model.Unit{},
		&model.Product{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}
	log.Println("Migrasi database berhasil.")
}

// SeedData mengisi data master awal
func SeedData(db *gorm.DB) {
	log.Println("Mengecek data seeder...")

	// Seeding Units (Satuan)
	var unitCount int64
	db.Model(&model.Unit{}).Count(&unitCount)
	if unitCount == 0 {
		log.Println("Menjalankan seeder untuk Units...")
		units := []model.Unit{
			{NamaUnit: "Pieces", Singkatan: "pcs"},
			{NamaUnit: "Kilogram", Singkatan: "kg"},
			{NamaUnit: "Box", Singkatan: "box"},
			{NamaUnit: "Liter", Singkatan: "ltr"},
		}
		// Hook BeforeCreate akan otomatis mengisi UUID
		if err := db.Create(&units).Error; err != nil {
			log.Fatal("Gagal seeding units:", err)
		}
	}

	// (Kita bisa tambahkan seeder Kategori di sini nanti)

	log.Println("Seeder selesai.")
}
