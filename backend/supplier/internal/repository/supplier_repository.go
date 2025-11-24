package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const cacheKeyAllSuppliers = "all_suppliers"

// Helper untuk membuat key cache per ID
func supplierCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("supplier:%s", id.String())
}

type SupplierRepository interface {
	Create(supplier *model.Supplier) error
	FindAll() ([]model.Supplier, error)
	FindByID(id uuid.UUID) (*model.Supplier, error)
	Update(supplier *model.Supplier) (*model.Supplier, error)
	Delete(id uuid.UUID) error
}

type supplierRepository struct {
	db       *gorm.DB
	rdb      *redis.Client
	ctx      context.Context
	cacheTtl time.Duration
}

func NewSupplierRepository(db *gorm.DB, rdb *redis.Client) SupplierRepository {
	return &supplierRepository{
		db:       db,
		rdb:      rdb,
		ctx:      context.Background(),
		cacheTtl: time.Minute * 30, // Cache kedaluwarsa dalam 30 menit
	}
}

func (r *supplierRepository) Create(supplier *model.Supplier) error {
	// 1. Simpan ke DB
	if err := r.db.Create(supplier).Error; err != nil {
		return err
	}

	// 2. Invalidate Cache
	log.Println("[CACHE DEBUG] Create Supplier: Menghapus cache list 'all_suppliers'")
	r.rdb.Del(r.ctx, cacheKeyAllSuppliers)

	return nil
}

func (r *supplierRepository) FindAll() ([]model.Supplier, error) {
	var suppliers []model.Supplier

	// 1. Cek Redis
	log.Println("[CACHE DEBUG] FindAll: Mencari data di Redis...")
	cached, err := r.rdb.Get(r.ctx, cacheKeyAllSuppliers).Result()

	if err == nil {
		log.Println("[CACHE DEBUG] FindAll: HIT! Mengembalikan data dari Redis.")
		json.Unmarshal([]byte(cached), &suppliers)
		return suppliers, nil
	}

	// 2. Cek DB (Cache Miss)
	log.Printf("[CACHE DEBUG] FindAll: MISS! (%v). Mengambil dari DB...", err)
	if err := r.db.Find(&suppliers).Error; err != nil {
		return nil, err
	}

	// 3. Simpan ke Redis
	data, _ := json.Marshal(suppliers)
	log.Println("[CACHE DEBUG] FindAll: Menyimpan data terbaru ke Redis.")
	r.rdb.Set(r.ctx, cacheKeyAllSuppliers, data, r.cacheTtl)

	return suppliers, nil
}

func (r *supplierRepository) FindByID(id uuid.UUID) (*model.Supplier, error) {
	var supplier model.Supplier
	key := supplierCacheKey(id)

	// 1. Cek Redis
	log.Printf("[CACHE DEBUG] FindByID: Mencari key '%s' di Redis...", key)
	cached, err := r.rdb.Get(r.ctx, key).Result()

	if err == nil {
		log.Printf("[CACHE DEBUG] FindByID: HIT! Key '%s' ditemukan.", key)
		json.Unmarshal([]byte(cached), &supplier)
		return &supplier, nil
	}

	// 2. Cek DB
	log.Printf("[CACHE DEBUG] FindByID: MISS! Mengambil ID %s dari DB...", id)
	if err := r.db.First(&supplier, id).Error; err != nil {
		return nil, err
	}

	// 3. Simpan ke Redis
	data, _ := json.Marshal(supplier)
	log.Printf("[CACHE DEBUG] FindByID: Menyimpan key '%s' ke Redis.", key)
	r.rdb.Set(r.ctx, key, data, r.cacheTtl)

	return &supplier, nil
}

func (r *supplierRepository) Update(supplier *model.Supplier) (*model.Supplier, error) {
	// 1. Update DB
	if err := r.db.Save(supplier).Error; err != nil {
		return nil, err
	}

	// 2. Invalidate Cache (Hapus List DAN Detail)
	key := supplierCacheKey(supplier.ID)
	log.Printf("[CACHE DEBUG] Update: Menghapus cache list '%s' dan detail '%s'", cacheKeyAllSuppliers, key)

	r.rdb.Del(r.ctx, cacheKeyAllSuppliers) // Hapus list agar data lama hilang
	r.rdb.Del(r.ctx, key)                  // Hapus detail agar data lama hilang

	return supplier, nil
}

func (r *supplierRepository) Delete(id uuid.UUID) error {
	// 1. Delete DB (Soft Delete)
	if err := r.db.Delete(&model.Supplier{}, id).Error; err != nil {
		return err
	}

	// 2. Invalidate Cache
	key := supplierCacheKey(id)
	log.Printf("[CACHE DEBUG] Delete: Menghapus cache list '%s' dan detail '%s'", cacheKeyAllSuppliers, key)

	r.rdb.Del(r.ctx, cacheKeyAllSuppliers)
	r.rdb.Del(r.ctx, key)

	return nil
}
