package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role merepresentasikan tabel 'roles'
type Role struct {
	ID        uint   `gorm:"primaryKey"`
	NamaPeran string `gorm:"unique;not null;size:50"`
	Users     []User // Relasi
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User merepresentasikan tabel 'users'
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Nama         string    `gorm:"not null;size:255"`
	Email        string    `gorm:"unique;not null;size:255"`
	PasswordHash string    `gorm:"not null;size:255"`
	RoleID       uint      `gorm:"not null"` // Foreign key ke Role
	Role         Role      `gorm:"foreignKey:RoleID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Hook GORM: Akan dipanggil sebelum Create
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID baru jika ID-nya kosong
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}
