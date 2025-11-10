package dto

import "github.com/google/uuid"

// --- Category DTOs ---
type CreateCategoryDTO struct {
	NamaKategori string `json:"nama_kategori" validate:"required,min=3"`
	Deskripsi    string `json:"deskripsi"`
}
type UpdateCategoryDTO struct {
	NamaKategori string `json:"nama_kategori" validate:"omitempty,min=3"`
	Deskripsi    string `json:"deskripsi"`
}

// --- Unit DTOs ---
type CreateUnitDTO struct {
	NamaUnit  string `json:"nama_unit" validate:"required"`
	Singkatan string `json:"singkatan" validate:"required,max=10"`
}
type UpdateUnitDTO struct {
	NamaUnit  string `json:"nama_unit" validate:"omitempty"`
	Singkatan string `json:"singkatan" validate:"omitempty,max=10"`
}

// --- Product DTOs ---
type CreateProductDTO struct {
	SKU        string    `json:"sku" validate:"required"`
	NamaProduk string    `json:"nama_produk" validate:"required,min=3"`
	Deskripsi  string    `json:"deskripsi"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	UnitID     uuid.UUID `json:"unit_id" validate:"required"`
}
type UpdateProductDTO struct {
	SKU        string    `json:"sku" validate:"omitempty"`
	NamaProduk string    `json:"nama_produk" validate:"omitempty,min=3"`
	Deskripsi  string    `json:"deskripsi"`
	CategoryID uuid.UUID `json:"category_id" validate:"omitempty"`
	UnitID     uuid.UUID `json:"unit_id" validate:"omitempty"`
}
