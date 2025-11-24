package service

import (
	"testing"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper setup
func setupTest() (SupplierService, *repository.SupplierRepositoryMock, *queue.EventPublisherMock) {
	mockRepo := new(repository.SupplierRepositoryMock)
	mockPub := new(queue.EventPublisherMock)
	service := NewSupplierService(mockRepo, mockPub)
	return service, mockRepo, mockPub
}

func TestCreate_Success(t *testing.T) {
	// 1. Arrange
	service, mockRepo, mockPub := setupTest()

	input := dto.CreateSupplierDTO{
		NamaPemasok: "PT Test",
		Email:       "test@pt.com",
		Telepon:     "08123",
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Supplier")).Return(nil)
	mockPub.On("Publish", "supplier.created", mock.Anything).Return(nil)

	// 2. Act
	res, err := service.Create(input)

	// --- FIX: Tunggu Goroutine ---
	time.Sleep(50 * time.Millisecond)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, input.NamaPemasok, res.NamaPemasok)

	mockRepo.AssertExpectations(t)
	// Assert bahwa Publish terpanggil
	mockPub.AssertCalled(t, "Publish", "supplier.created", mock.Anything)
}

func TestFindAll_Success(t *testing.T) {
	service, mockRepo, _ := setupTest()

	mockData := []model.Supplier{
		{NamaPemasok: "A"},
		{NamaPemasok: "B"},
	}

	mockRepo.On("FindAll").Return(mockData, nil)

	res, err := service.FindAll()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestFindByID_Success(t *testing.T) {
	service, mockRepo, _ := setupTest()
	id := uuid.New()
	mockData := &model.Supplier{ID: id, NamaPemasok: "A"}

	mockRepo.On("FindByID", id).Return(mockData, nil)

	res, err := service.FindByID(id)

	assert.NoError(t, err)
	assert.Equal(t, id, res.ID)
}

func TestUpdate_Success(t *testing.T) {
	// 1. Arrange
	service, mockRepo, mockPub := setupTest()
	id := uuid.New()

	input := dto.UpdateSupplierDTO{
		NamaPemasok: "PT Updated",
	}

	existingData := &model.Supplier{ID: id, NamaPemasok: "PT Lama"}
	updatedData := &model.Supplier{ID: id, NamaPemasok: "PT Updated"}

	mockRepo.On("FindByID", id).Return(existingData, nil)
	mockRepo.On("Update", existingData).Return(updatedData, nil)
	mockPub.On("Publish", "supplier.updated", updatedData).Return(nil)

	// 2. Act
	res, err := service.Update(id, input)

	// --- FIX: Tunggu Goroutine ---
	time.Sleep(50 * time.Millisecond)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, "PT Updated", res.NamaPemasok)
	mockPub.AssertCalled(t, "Publish", "supplier.updated", updatedData)
}

func TestDelete_Success(t *testing.T) {
	// 1. Arrange
	service, mockRepo, mockPub := setupTest()
	id := uuid.New()
	existingData := &model.Supplier{ID: id}

	mockRepo.On("FindByID", id).Return(existingData, nil)
	mockRepo.On("Delete", id).Return(nil)
	mockPub.On("Publish", "supplier.deleted", existingData).Return(nil)

	// 2. Act
	err := service.Delete(id)

	// --- FIX: Tunggu Goroutine ---
	time.Sleep(50 * time.Millisecond)

	// 3. Assert
	assert.NoError(t, err)
	mockPub.AssertCalled(t, "Publish", "supplier.deleted", existingData)
}
