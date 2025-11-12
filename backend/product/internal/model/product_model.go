package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category merepresentasikan tabel 'categories'
type Category struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	NamaKategori string    `gorm:"unique;not null;size:100"`
	Deskripsi    string
	Products     []Product // Relasi one-to-many
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Unit merepresentasikan tabel 'units' (satuan)
type Unit struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	NamaUnit  string    `gorm:"unique;not null;size:50"` // misal: "Kilogram"
	Singkatan string    `gorm:"unique;not null;size:10"` // misal: "Kg"
	Products  []Product // Relasi one-to-many
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Product merepresentasikan tabel 'products' (data master)
type Product struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	SKU        string    `gorm:"unique;not null;size:100"`
	NamaProduk string    `gorm:"not null;size:255"`
	Deskripsi  string
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	UnitID     uuid.UUID `gorm:"type:uuid;not null"`
	Category   Category  `gorm:"foreignKey:CategoryID"` // Relasi many-to-one
	Unit       Unit      `gorm:"foreignKey:UnitID"`     // Relasi many-to-one
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// --- Hook GORM untuk auto-generate UUID ---
func (model *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	return
}

func (model *Unit) BeforeCreate(tx *gorm.DB) (err error) {
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	return
}

func (model *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	return
}
