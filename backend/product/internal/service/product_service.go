package service

import (
	"errors"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Definisikan interface
type ProductService interface {
	// Category
	CreateCategory(input dto.CreateCategoryDTO) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(id uuid.UUID, input dto.UpdateCategoryDTO) (*model.Category, error)
	DeleteCategory(id uuid.UUID) error

	// Unit
	CreateUnit(input dto.CreateUnitDTO) (*model.Unit, error)
	GetAllUnits() ([]model.Unit, error)
	UpdateUnit(id uuid.UUID, input dto.UpdateUnitDTO) (*model.Unit, error)
	DeleteUnit(id uuid.UUID) error

	// Product
	CreateProduct(input dto.CreateProductDTO) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id uuid.UUID) (*model.Product, error)
	UpdateProduct(id uuid.UUID, input dto.UpdateProductDTO) (*model.Product, error)
	DeleteProduct(id uuid.UUID) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// --- Implementasi Category ---
func (s *productService) CreateCategory(input dto.CreateCategoryDTO) (*model.Category, error) {
	newCategory := &model.Category{
		NamaKategori: input.NamaKategori,
		Deskripsi:    input.Deskripsi,
	}
	err := s.repo.CreateCategory(newCategory)
	return newCategory, err
}

func (s *productService) GetAllCategories() ([]model.Category, error) {
	return s.repo.FindAllCategories()
}

func (s *productService) UpdateCategory(id uuid.UUID, input dto.UpdateCategoryDTO) (*model.Category, error) {
	// Ambil data
	category, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err // Akan error jika not found
	}

	// Perbarui field jika ada di input
	if input.NamaKategori != "" {
		category.NamaKategori = input.NamaKategori
	}
	if input.Deskripsi != "" {
		category.Deskripsi = input.Deskripsi
	}

	// Simpan (Repo akan invalidate cache)
	return s.repo.UpdateCategory(category)
}

func (s *productService) DeleteCategory(id uuid.UUID) error {
	// Cek apakah ada
	_, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err // Error not found
	}
	// Hapus (Repo akan invalidate cache)
	return s.repo.DeleteCategory(id)
}

// --- Implementasi Unit ---
func (s *productService) CreateUnit(input dto.CreateUnitDTO) (*model.Unit, error) {
	newUnit := &model.Unit{
		NamaUnit:  input.NamaUnit,
		Singkatan: input.Singkatan,
	}
	err := s.repo.CreateUnit(newUnit)
	return newUnit, err
}

func (s *productService) GetAllUnits() ([]model.Unit, error) {
	return s.repo.FindAllUnits()
}

func (s *productService) UpdateUnit(id uuid.UUID, input dto.UpdateUnitDTO) (*model.Unit, error) {
	unit, err := s.repo.FindUnitByID(id) // Perlu FindUnitByID di repo
	if err != nil {
		return nil, err
	}
	if input.NamaUnit != "" {
		unit.NamaUnit = input.NamaUnit
	}
	if input.Singkatan != "" {
		unit.Singkatan = input.Singkatan
	}
	return s.repo.UpdateUnit(unit)
}

func (s *productService) DeleteUnit(id uuid.UUID) error {
	_, err := s.repo.FindUnitByID(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteUnit(id)
}

// --- Implementasi Product ---
func (s *productService) CreateProduct(input dto.CreateProductDTO) (*model.Product, error) {
	newProduct := &model.Product{
		SKU:        input.SKU,
		NamaProduk: input.NamaProduk,
		Deskripsi:  input.Deskripsi,
		CategoryID: input.CategoryID,
		UnitID:     input.UnitID,
	}

	err := s.repo.CreateProduct(newProduct)
	return newProduct, err
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAllProducts()
}

func (s *productService) GetProductByID(id uuid.UUID) (*model.Product, error) {
	return s.repo.FindProductByID(id)
}

func (s *productService) UpdateProduct(id uuid.UUID, input dto.UpdateProductDTO) (*model.Product, error) {
	product, err := s.repo.FindProductByID(id)
	if err != nil {
		return nil, err
	}
	if input.SKU != "" {
		product.SKU = input.SKU
	}
	if input.NamaProduk != "" {
		product.NamaProduk = input.NamaProduk
	}
	if input.Deskripsi != "" {
		product.Deskripsi = input.Deskripsi
	}
	if input.CategoryID != uuid.Nil {
		product.CategoryID = input.CategoryID
	}
	if input.UnitID != uuid.Nil {
		product.UnitID = input.UnitID
	}
	return s.repo.UpdateProduct(product)
}

func (s *productService) DeleteProduct(id uuid.UUID) error {
	_, err := s.repo.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Produk tidak ditemukan")
		}
		return err
	}
	return s.repo.DeleteProduct(id)
}
