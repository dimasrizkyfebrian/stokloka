package service

import (
	"errors"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductService Interface
type ProductService interface {
	// Category
	CreateCategory(input dto.CreateCategoryDTO) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(id uuid.UUID, input dto.UpdateCategoryDTO) (*model.Category, error)
	DeleteCategory(id uuid.UUID) error

	// Unit
	CreateUnit(input dto.CreateUnitDTO) (*model.Unit, error)
	GetAllUnits() ([]model.Unit, error)
	FindUnitByID(id uuid.UUID) (*model.Unit, error)
	UpdateUnit(id uuid.UUID, input dto.UpdateUnitDTO) (*model.Unit, error)
	DeleteUnit(id uuid.UUID) error

	// Product
	CreateProduct(input dto.CreateProductDTO) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id uuid.UUID) (*model.Product, error)
	UpdateProduct(id uuid.UUID, input dto.UpdateProductDTO) (*model.Product, error)
	DeleteProduct(id uuid.UUID) error
}

// productService Struct
type productService struct {
	repo      repository.ProductRepository
	publisher queue.EventPublisher
}

// NewProductService "Constructor"
func NewProductService(repo repository.ProductRepository, publisher queue.EventPublisher) ProductService {
	return &productService{
		repo:      repo,
		publisher: publisher,
	}
}

// --- Implementasi Category ---
func (s *productService) CreateCategory(input dto.CreateCategoryDTO) (*model.Category, error) {
	newCategory := &model.Category{
		NamaKategori: input.NamaKategori,
		Deskripsi:    input.Deskripsi,
	}
	err := s.repo.CreateCategory(newCategory)
	if err != nil {
		return nil, err
	}
	// Panggil publisher
	if err := s.publisher.Publish("category.created", newCategory); err != nil {
		log.Printf("Gagal publish event category.created: %v", err)
	}
	return newCategory, err
}

func (s *productService) GetAllCategories() ([]model.Category, error) {
	return s.repo.FindAllCategories()
}

func (s *productService) UpdateCategory(id uuid.UUID, input dto.UpdateCategoryDTO) (*model.Category, error) {
	category, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if input.NamaKategori != "" {
		category.NamaKategori = input.NamaKategori
	}
	if input.Deskripsi != "" {
		category.Deskripsi = input.Deskripsi
	}

	updatedCategory, err := s.repo.UpdateCategory(category)
	if err != nil {
		return nil, err
	}
	// Panggil publisher
	if err := s.publisher.Publish("category.updated", updatedCategory); err != nil {
		log.Printf("Gagal publish event category.updated: %v", err)
	}
	return updatedCategory, nil
}

func (s *productService) DeleteCategory(id uuid.UUID) error {
	category, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteCategory(id)
	if err != nil {
		return err
	}
	// Panggil publisher
	if err := s.publisher.Publish("category.deleted", category); err != nil {
		log.Printf("Gagal publish event category.deleted: %v", err)
	}
	return nil
}

// --- Implementasi Unit ---
func (s *productService) CreateUnit(input dto.CreateUnitDTO) (*model.Unit, error) {
	newUnit := &model.Unit{
		NamaUnit:  input.NamaUnit,
		Singkatan: input.Singkatan,
	}
	err := s.repo.CreateUnit(newUnit)
	if err != nil {
		return nil, err
	}
	// Panggil publisher
	if err := s.publisher.Publish("unit.deleted", newUnit); err != nil {
		log.Printf("Gagal publish event unit.deleted: %v", err)
	}
	return newUnit, err
}

func (s *productService) GetAllUnits() ([]model.Unit, error) {
	return s.repo.FindAllUnits()
}

func (s *productService) FindUnitByID(id uuid.UUID) (*model.Unit, error) {
	return s.repo.FindUnitByID(id)
}

func (s *productService) UpdateUnit(id uuid.UUID, input dto.UpdateUnitDTO) (*model.Unit, error) {
	unit, err := s.repo.FindUnitByID(id)
	if err != nil {
		return nil, err
	}
	if input.NamaUnit != "" {
		unit.NamaUnit = input.NamaUnit
	}
	if input.Singkatan != "" {
		unit.Singkatan = input.Singkatan
	}

	updatedUnit, err := s.repo.UpdateUnit(unit)
	if err != nil {
		return nil, err
	}
	// Panggil publisher
	if err := s.publisher.Publish("unit.deleted", unit); err != nil {
		log.Printf("Gagal publish event unit.deleted: %v", err)
	}
	return updatedUnit, nil
}

func (s *productService) DeleteUnit(id uuid.UUID) error {
	unit, err := s.repo.FindUnitByID(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteUnit(id)
	if err != nil {
		return err
	}
	// Panggil publisher
	if err := s.publisher.Publish("unit.deleted", unit); err != nil {
		log.Printf("Gagal publish event unit.deleted: %v", err)
	}
	return nil
}

// --- Implementasi Product ---
func (s *productService) CreateProduct(input dto.CreateProductDTO) (*model.Product, error) {
	// Cek dependensi
	if _, err := s.repo.FindCategoryByID(input.CategoryID); err != nil {
		return nil, errors.New("Category ID tidak valid")
	}
	if _, err := s.repo.FindUnitByID(input.UnitID); err != nil {
		return nil, errors.New("Unit ID tidak valid")
	}

	newProduct := &model.Product{
		SKU:        input.SKU,
		NamaProduk: input.NamaProduk,
		Deskripsi:  input.Deskripsi,
		CategoryID: input.CategoryID,
		UnitID:     input.UnitID,
	}

	err := s.repo.CreateProduct(newProduct)
	if err != nil {
		return nil, err
	}

	createdProduct, err := s.repo.FindProductByID(newProduct.ID) // Preload
	// Panggil publisher
	if err != nil {
		log.Printf("Gagal preload produk untuk event: %v", err)
		s.publisher.Publish("product.created", newProduct)
	} else {
		s.publisher.Publish("product.created", createdProduct)
	}

	return createdProduct, nil
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
		if _, err := s.repo.FindCategoryByID(input.CategoryID); err != nil {
			return nil, errors.New("Category ID tidak valid")
		}
		product.CategoryID = input.CategoryID
	}
	if input.UnitID != uuid.Nil {
		if _, err := s.repo.FindUnitByID(input.UnitID); err != nil {
			return nil, errors.New("Unit ID tidak valid")
		}
		product.UnitID = input.UnitID
	}

	updatedProduct, err := s.repo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	finalProduct, err := s.repo.FindProductByID(updatedProduct.ID) // Preload
	// Panggil publisher
	if err != nil {
		log.Printf("Gagal preload produk untuk event: %v", err)
		s.publisher.Publish("product.updated", updatedProduct)
	} else {
		s.publisher.Publish("product.updated", finalProduct)
	}

	return finalProduct, nil
}

func (s *productService) DeleteProduct(id uuid.UUID) error {
	product, err := s.repo.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Produk tidak ditemukan")
		}
		return err
	}

	err = s.repo.DeleteProduct(id)
	if err != nil {
		return err
	}
	// Panggil publisher
	if err := s.publisher.Publish("product.deleted", product); err != nil {
		log.Printf("Gagal publish event product.deleted: %v", err)
	}
	return nil
}
