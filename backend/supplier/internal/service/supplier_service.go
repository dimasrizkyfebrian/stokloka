package service

import (
	"errors"
	"log"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupplierService interface {
	Create(input dto.CreateSupplierDTO) (*model.Supplier, error)
	FindAll() ([]model.Supplier, error)
	FindByID(id uuid.UUID) (*model.Supplier, error)
	Update(id uuid.UUID, input dto.UpdateSupplierDTO) (*model.Supplier, error)
	Delete(id uuid.UUID) error
}

type supplierService struct {
	repo      repository.SupplierRepository
	publisher queue.EventPublisher
}

func NewSupplierService(repo repository.SupplierRepository, pub queue.EventPublisher) SupplierService {
	return &supplierService{repo: repo, publisher: pub}
}

// Helper Goroutine Publish
func (s *supplierService) publishAsync(key string, data interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic di publisher: %v", r)
			}
		}()
		if err := s.publisher.Publish(key, data); err != nil {
			log.Printf("Gagal publish %s: %v", key, err)
		}
	}()
}

func (s *supplierService) Create(input dto.CreateSupplierDTO) (*model.Supplier, error) {
	supplier := &model.Supplier{
		KodePemasok:  input.KodePemasok,
		NamaPemasok:  input.NamaPemasok,
		KontakPerson: input.KontakPerson,
		Email:        input.Email,
		Telepon:      input.Telepon,
		Alamat:       input.Alamat,
	}

	if err := s.repo.Create(supplier); err != nil {
		return nil, err
	}

	// Publish Event (Async)
	s.publishAsync("supplier.created", supplier)

	return supplier, nil
}

func (s *supplierService) FindAll() ([]model.Supplier, error) {
	return s.repo.FindAll()
}

func (s *supplierService) FindByID(id uuid.UUID) (*model.Supplier, error) {
	return s.repo.FindByID(id)
}

func (s *supplierService) Update(id uuid.UUID, input dto.UpdateSupplierDTO) (*model.Supplier, error) {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if present
	if input.KodePemasok != "" {
		supplier.KodePemasok = input.KodePemasok
	}
	if input.NamaPemasok != "" {
		supplier.NamaPemasok = input.NamaPemasok
	}
	if input.KontakPerson != "" {
		supplier.KontakPerson = input.KontakPerson
	}
	if input.Email != "" {
		supplier.Email = input.Email
	}
	if input.Telepon != "" {
		supplier.Telepon = input.Telepon
	}
	if input.Alamat != "" {
		supplier.Alamat = input.Alamat
	}

	updated, err := s.repo.Update(supplier)
	if err != nil {
		return nil, err
	}

	s.publishAsync("supplier.updated", updated)
	return updated, nil
}

func (s *supplierService) Delete(id uuid.UUID) error {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("supplier tidak ditemukan")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.publishAsync("supplier.deleted", supplier)
	return nil
}
