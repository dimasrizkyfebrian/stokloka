package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Definisikan cache keys
const (
	cacheKeyAllCategories = "all_categories"
	cacheKeyAllUnits      = "all_units"
	cacheKeyAllProducts   = "all_products"
)

// Helper untuk cache key produk individual
func categoryCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("category:%s", id.String())
}
func unitCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("unit:%s", id.String())
}
func productCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("product:%s", id.String())
}

// Definisikan interface
type ProductRepository interface {
	// Category
	CreateCategory(category *model.Category) error
	FindCategoryByID(id uuid.UUID) (*model.Category, error)
	FindAllCategories() ([]model.Category, error)
	UpdateCategory(category *model.Category) (*model.Category, error)
	DeleteCategory(id uuid.UUID) error

	// Unit
	CreateUnit(unit *model.Unit) error
	FindUnitByID(id uuid.UUID) (*model.Unit, error)
	FindAllUnits() ([]model.Unit, error)
	UpdateUnit(unit *model.Unit) (*model.Unit, error)
	DeleteUnit(id uuid.UUID) error

	// Product
	CreateProduct(product *model.Product) error
	FindProductByID(id uuid.UUID) (*model.Product, error)
	FindAllProducts() ([]model.Product, error)
	UpdateProduct(product *model.Product) (*model.Product, error)
	DeleteProduct(id uuid.UUID) error
}

// Struct implementasi
type productRepository struct {
	db       *gorm.DB      // Koneksi Postgres
	rdb      *redis.Client // Koneksi Redis
	ctx      context.Context
	cacheTtl time.Duration // Waktu hidup cache
}

// "Constructor"
func NewProductRepository(db *gorm.DB, rdb *redis.Client) ProductRepository {
	return &productRepository{
		db:       db,
		rdb:      rdb,
		ctx:      context.Background(),
		cacheTtl: time.Minute * 10, // Simpan cache selama 10 menit
	}
}

// --- Implementasi Category ---
func (r *productRepository) CreateCategory(category *model.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	log.Println("CACHE: Menghapus cache", cacheKeyAllCategories)
	r.rdb.Del(r.ctx, cacheKeyAllCategories) // Invalidate cache
	return nil
}

func (r *productRepository) FindCategoryByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	cacheKey := categoryCacheKey(id)

	cachedData, err := r.rdb.Get(r.ctx, cacheKey).Result()
	if err == nil {
		log.Println("CACHE: HIT! Mengembalikan data category (by ID) dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &category); err != nil {
			return nil, err
		}
		return &category, nil
	}

	log.Println("CACHE: MISS! Mengambil data category (by ID) dari PostgreSQL...")
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(category)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKey, dataToCache, r.cacheTtl)
	}
	return &category, nil
}

func (r *productRepository) FindAllCategories() ([]model.Category, error) {
	var categories []model.Category
	cacheKey := cacheKeyAllCategories

	log.Println("CACHE: Mencari data dari Redis...")
	cachedData, err := r.rdb.Get(r.ctx, cacheKey).Result()

	if err == nil {
		log.Println("CACHE: HIT! Mengembalikan data dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &categories); err != nil {
			return nil, err
		}
		return categories, nil
	}

	log.Println("CACHE: MISS! Mengambil data dari PostgreSQL...")

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(categories)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKey, dataToCache, r.cacheTtl)
	}

	return categories, nil
}

func (r *productRepository) UpdateCategory(category *model.Category) (*model.Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	// Invalidate cache
	log.Println("CACHE: Menghapus cache", cacheKeyAllCategories)
	r.rdb.Del(r.ctx, cacheKeyAllCategories)
	log.Println("CACHE: Menghapus cache", categoryCacheKey(category.ID))
	r.rdb.Del(r.ctx, categoryCacheKey(category.ID))
	return category, nil
}

func (r *productRepository) DeleteCategory(id uuid.UUID) error {
	// Invalidate cache SEBELUM delete
	log.Println("CACHE: Menghapus cache", cacheKeyAllCategories)
	r.rdb.Del(r.ctx, cacheKeyAllCategories)
	log.Println("CACHE: Menghapus cache", categoryCacheKey(id))
	r.rdb.Del(r.ctx, categoryCacheKey(id))

	// Hard delete (karena model Category tidak punya gorm.DeletedAt)
	return r.db.Delete(&model.Category{}, id).Error
}

// --- Implementasi Unit ---
func (r *productRepository) CreateUnit(unit *model.Unit) error {
	if err := r.db.Create(unit).Error; err != nil {
		return err
	}
	// Invalidate cache "all units"
	log.Println("CACHE: Menghapus cache", cacheKeyAllUnits)
	r.rdb.Del(r.ctx, cacheKeyAllUnits)
	return nil
}

func (r *productRepository) FindAllUnits() ([]model.Unit, error) {
	var units []model.Unit

	// Coba ambil dari REDIS
	cachedData, err := r.rdb.Get(r.ctx, cacheKeyAllUnits).Result()
	if err == nil {
		// CACHE HIT
		log.Println("CACHE: HIT! Mengembalikan data units dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &units); err != nil {
			return nil, err
		}
		return units, nil
	}

	// CACHE MISS -> Ambil dari POSTGRES
	log.Println("CACHE: MISS! Mengambil data units dari PostgreSQL...")
	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}

	// Simpan ke REDIS
	dataToCache, err := json.Marshal(units)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKeyAllUnits, dataToCache, r.cacheTtl)
	}
	return units, nil
}

func (r *productRepository) FindUnitByID(id uuid.UUID) (*model.Unit, error) {
	var unit model.Unit
	cacheKey := unitCacheKey(id)

	cachedData, err := r.rdb.Get(r.ctx, cacheKey).Result()
	if err == nil {
		log.Println("CACHE: HIT! Mengembalikan data unit (by ID) dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &unit); err != nil {
			return nil, err
		}
		return &unit, nil
	}

	log.Println("CACHE: MISS! Mengambil data unit (by ID) dari PostgreSQL...")
	if err := r.db.First(&unit, id).Error; err != nil {
		return nil, err
	}

	dataToCache, err := json.Marshal(unit)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKey, dataToCache, r.cacheTtl)
	}
	return &unit, nil
}

func (r *productRepository) UpdateUnit(unit *model.Unit) (*model.Unit, error) {
	if err := r.db.Save(unit).Error; err != nil {
		return nil, err
	}
	// Invalidate cache
	log.Println("CACHE: Menghapus cache", cacheKeyAllUnits)
	r.rdb.Del(r.ctx, cacheKeyAllUnits)
	log.Println("CACHE: Menghapus cache", unitCacheKey(unit.ID))
	r.rdb.Del(r.ctx, unitCacheKey(unit.ID))
	return unit, nil
}

func (r *productRepository) DeleteUnit(id uuid.UUID) error {
	// Invalidate cache
	log.Println("CACHE: Menghapus cache", cacheKeyAllUnits)
	r.rdb.Del(r.ctx, cacheKeyAllUnits)
	log.Println("CACHE: Menghapus cache", unitCacheKey(id))
	r.rdb.Del(r.ctx, unitCacheKey(id))

	// Hard delete
	return r.db.Delete(&model.Unit{}, id).Error
}

// --- Implementasi Product ---
func (r *productRepository) CreateProduct(product *model.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}
	// Invalidate cache "all products"
	log.Println("CACHE: Menghapus cache", cacheKeyAllProducts)
	r.rdb.Del(r.ctx, cacheKeyAllProducts)
	return nil
}

func (r *productRepository) FindAllProducts() ([]model.Product, error) {
	var products []model.Product

	// Coba ambil dari REDIS
	cachedData, err := r.rdb.Get(r.ctx, cacheKeyAllProducts).Result()
	if err == nil {
		// CACHE HIT
		log.Println("CACHE: HIT! Mengembalikan data products dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &products); err != nil {
			return nil, err
		}
		return products, nil
	}

	// CACHE MISS -> Ambil dari POSTGRES
	log.Println("CACHE: MISS! Mengambil data products dari PostgreSQL...")
	// Preload relasinya (Category dan Unit) agar datanya lengkap
	if err := r.db.Preload("Category").Preload("Unit").Find(&products).Error; err != nil {
		return nil, err
	}

	// Simpan ke REDIS
	dataToCache, err := json.Marshal(products)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKeyAllProducts, dataToCache, r.cacheTtl)
	}
	return products, nil
}

func (r *productRepository) FindProductByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	cacheKey := productCacheKey(id) // Buat cache key unik, cth: "product:uuid-..."

	// Coba ambil dari REDIS
	cachedData, err := r.rdb.Get(r.ctx, cacheKey).Result()
	if err == nil {
		// CACHE HIT
		log.Println("CACHE: HIT! Mengembalikan data product (by ID) dari Redis.")
		if err := json.Unmarshal([]byte(cachedData), &product); err != nil {
			return nil, err
		}
		return &product, nil
	}

	// CACHE MISS -> Ambil dari POSTGRES
	log.Println("CACHE: MISS! Mengambil data product (by ID) dari PostgreSQL...")
	if err := r.db.Preload("Category").Preload("Unit").First(&product, id).Error; err != nil {
		return nil, err
	}

	// Simpan ke REDIS
	dataToCache, err := json.Marshal(product)
	if err == nil {
		r.rdb.Set(r.ctx, cacheKey, dataToCache, r.cacheTtl)
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(product *model.Product) (*model.Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		return nil, err
	}
	// Invalidate cache
	log.Println("CACHE: Menghapus cache", cacheKeyAllProducts)
	r.rdb.Del(r.ctx, cacheKeyAllProducts)
	log.Println("CACHE: Menghapus cache", productCacheKey(product.ID))
	r.rdb.Del(r.ctx, productCacheKey(product.ID))
	return product, nil
}

func (r *productRepository) DeleteProduct(id uuid.UUID) error {
	// Invalidate cache
	log.Println("CACHE: Menghapus cache", cacheKeyAllProducts)
	r.rdb.Del(r.ctx, cacheKeyAllProducts)
	log.Println("CACHE: Menghapus cache", productCacheKey(id))
	r.rdb.Del(r.ctx, productCacheKey(id))

	// Soft delete (karena model Product punya gorm.DeletedAt)
	return r.db.Delete(&model.Product{}, id).Error
}
