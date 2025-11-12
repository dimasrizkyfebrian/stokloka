package repository

import (
	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

// Implementasi MOCK untuk SEMUA metode di interface
func (m *ProductRepositoryMock) CreateCategory(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}
func (m *ProductRepositoryMock) FindCategoryByID(id uuid.UUID) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}
func (m *ProductRepositoryMock) FindAllCategories() ([]model.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Category), args.Error(1)
}
func (m *ProductRepositoryMock) UpdateCategory(category *model.Category) (*model.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}
func (m *ProductRepositoryMock) DeleteCategory(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) CreateUnit(unit *model.Unit) error {
	args := m.Called(unit)
	return args.Error(0)
}
func (m *ProductRepositoryMock) FindUnitByID(id uuid.UUID) (*model.Unit, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Unit), args.Error(1)
}
func (m *ProductRepositoryMock) FindAllUnits() ([]model.Unit, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Unit), args.Error(1)
}
func (m *ProductRepositoryMock) UpdateUnit(unit *model.Unit) (*model.Unit, error) {
	args := m.Called(unit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Unit), args.Error(1)
}
func (m *ProductRepositoryMock) DeleteUnit(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) CreateProduct(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *ProductRepositoryMock) FindProductByID(id uuid.UUID) (*model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}
func (m *ProductRepositoryMock) FindAllProducts() ([]model.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}
func (m *ProductRepositoryMock) UpdateProduct(product *model.Product) (*model.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}
func (m *ProductRepositoryMock) DeleteProduct(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
