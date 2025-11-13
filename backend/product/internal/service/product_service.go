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

// ProductService Interface (Lengkap)
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
	publisher queue.EventPublisher // Menggunakan Interface
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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Category) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish category.created: %v", r)
			}
		}()
		if err := s.publisher.Publish("category.created", data); err != nil {
			log.Printf("Gagal publish event category.created (background): %v", err)
		}
	}(*newCategory) // Kirim copy-an data

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Category) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish category.updated: %v", r)
			}
		}()
		if err := s.publisher.Publish("category.updated", data); err != nil {
			log.Printf("Gagal publish event category.updated (background): %v", err)
		}
	}(*updatedCategory)

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Category) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish category.deleted: %v", r)
			}
		}()
		if err := s.publisher.Publish("category.deleted", data); err != nil {
			log.Printf("Gagal publish event category.deleted (background): %v", err)
		}
	}(*category)

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Unit) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish unit.created: %v", r)
			}
		}()
		if err := s.publisher.Publish("unit.created", data); err != nil {
			log.Printf("Gagal publish event unit.created (background): %v", err)
		}
	}(*newUnit)

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Unit) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish unit.updated: %v", r)
			}
		}()
		if err := s.publisher.Publish("unit.updated", data); err != nil {
			log.Printf("Gagal publish event unit.updated (background): %v", err)
		}
	}(*updatedUnit)

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Unit) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish unit.deleted: %v", r)
			}
		}()
		if err := s.publisher.Publish("unit.deleted", data); err != nil {
			log.Printf("Gagal publish event unit.deleted (background): %v", err)
		}
	}(*unit)

	return nil
}

// --- Implementasi Product ---
func (s *productService) CreateProduct(input dto.CreateProductDTO) (*model.Product, error) {
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

	createdProduct, err := s.repo.FindProductByID(newProduct.ID)
	if err != nil {
		log.Printf("Gagal preload produk untuk event: %v", err)
		// Tetap publish data seadanya
		go func(data model.Product) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic terdeteksi di goroutine publish product.created: %v", r)
				}
			}()
			if err := s.publisher.Publish("product.created", data); err != nil {
				log.Printf("Gagal publish event product.created (background): %v", err)
			}
		}(*newProduct)
		return newProduct, nil // Kembalikan data yg seadanya
	}

	// --- GOROUTINE PUBLISH (dengan data preload) ---
	go func(data model.Product) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish product.created: %v", r)
			}
		}()
		if err := s.publisher.Publish("product.created", data); err != nil {
			log.Printf("Gagal publish event product.created (background): %v", err)
		}
	}(*createdProduct)

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

	finalProduct, err := s.repo.FindProductByID(updatedProduct.ID)
	if err != nil {
		log.Printf("Gagal preload produk untuk event: %v", err)
		// --- GOROUTINE PUBLISH (data seadanya) ---
		go func(data model.Product) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic terdeteksi di goroutine publish product.updated: %v", r)
				}
			}()
			if err := s.publisher.Publish("product.updated", data); err != nil {
				log.Printf("Gagal publish event product.updated (background): %v", err)
			}
		}(*updatedProduct)
		return updatedProduct, nil
	}

	// --- GOROUTINE PUBLISH (dengan data preload) ---
	go func(data model.Product) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish product.updated: %v", r)
			}
		}()
		if err := s.publisher.Publish("product.updated", data); err != nil {
			log.Printf("Gagal publish event product.updated (background): %v", err)
		}
	}(*finalProduct)

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

	// --- GOROUTINE PUBLISH ---
	go func(data model.Product) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic terdeteksi di goroutine publish product.deleted: %v", r)
			}
		}()
		if err := s.publisher.Publish("product.deleted", data); err != nil {
			log.Printf("Gagal publish event product.deleted (background): %v", err)
		}
	}(*product)

	return nil
}
