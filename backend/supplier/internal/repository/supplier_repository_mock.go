package repository

import (
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SupplierRepositoryMock struct {
	mock.Mock
}

func (m *SupplierRepositoryMock) Create(supplier *model.Supplier) error {
	args := m.Called(supplier)
	return args.Error(0)
}

func (m *SupplierRepositoryMock) FindAll() ([]model.Supplier, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Supplier), args.Error(1)
}

func (m *SupplierRepositoryMock) FindByID(id uuid.UUID) (*model.Supplier, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Supplier), args.Error(1)
}

func (m *SupplierRepositoryMock) Update(supplier *model.Supplier) (*model.Supplier, error) {
	args := m.Called(supplier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Supplier), args.Error(1)
}

func (m *SupplierRepositoryMock) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
