package database

import (
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase membuat koneksi, migrasi, dan seeder
func InitDatabase(cfg *config.Config) *gorm.DB {
	// 1. DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	// 2. Buka Koneksi GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log semua query SQL
	})

	if err != nil {
		log.Fatalf("Gagal konek ke database %s: %v", cfg.DBName, err)
	}
	log.Printf("Koneksi ke database %s berhasil.", cfg.DBName)

	// 3. Jalankan Migrasi
	runMigrations(db)

	// 4. Jalankan Seeder
	SeedData(db, cfg)

	return db
}

// runMigrations menjalankan GORM AutoMigrate
func runMigrations(db *gorm.DB) {
	log.Println("Menjalankan migrasi database...")
	err := db.AutoMigrate(
		&model.Role{},
		&model.User{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}
	log.Println("Migrasi database berhasil.")
}

// SeedData mengisi data awal (Roles & Admin)
func SeedData(db *gorm.DB, cfg *config.Config) {
	log.Println("Mengecek data seeder...")

	// 1. Seeding Roles
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	if roleCount == 0 {
		log.Println("Menjalankan seeder untuk Roles...")
		roles := []model.Role{
			{ID: 1, NamaPeran: "admin"},
			{ID: 2, NamaPeran: "manager"},
			{ID: 3, NamaPeran: "staff"},
		}
		if err := db.Create(&roles).Error; err != nil {
			log.Fatal("Gagal seeding roles:", err)
		}
	}

	// 2. Seeding Admin Pertama
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		log.Println("Menjalankan seeder untuk Admin...")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)

		adminUser := model.User{
			Nama:         "Admin StokLoka",
			Email:        cfg.AdminEmail,
			PasswordHash: string(hashedPassword),
			RoleID:       1, // ID Role "admin"
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatal("Gagal seeding admin:", err)
		}
	}
	log.Println("Seeder selesai.")
}
