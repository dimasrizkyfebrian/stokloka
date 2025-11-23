package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Supplier struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	KodePemasok  string    `gorm:"unique;size:50"`
	NamaPemasok  string    `gorm:"not null;size:255"`
	KontakPerson string    `gorm:"size:100"`
	Email        string    `gorm:"size:255"`
	Telepon      string    `gorm:"size:50"`
	Alamat       string    `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Hook GORM untuk auto-generate UUID
func (s *Supplier) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
