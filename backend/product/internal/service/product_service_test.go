package service

import (
	"testing"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Helper setup
func setupProductServiceTest() (ProductService, *repository.ProductRepositoryMock, *queue.EventPublisherMock) {
	mockRepo := new(repository.ProductRepositoryMock)
	mockPublisher := new(queue.EventPublisherMock)
	authService := NewProductService(mockRepo, mockPublisher)
	return authService, mockRepo, mockPublisher
}

// --- Tes Produk ---
func TestCreateProduct_Success(t *testing.T) {
	// Arrange
	service, mockRepo, mockPublisher := setupProductServiceTest()

	inputDTO := dto.CreateProductDTO{
		SKU:        "TEST-001",
		NamaProduk: "Produk Tes",
		CategoryID: uuid.New(),
		UnitID:     uuid.New(),
	}

	// Buat data palsu untuk di-return
	mockCategory := &model.Category{ID: inputDTO.CategoryID}
	mockUnit := &model.Unit{ID: inputDTO.UnitID}

	// "Program" mock-nya
	mockRepo.On("FindCategoryByID", inputDTO.CategoryID).Return(mockCategory, nil)
	mockRepo.On("FindUnitByID", inputDTO.UnitID).Return(mockUnit, nil)
	mockRepo.On("CreateProduct", mock.AnythingOfType("*model.Product")).Return(nil)

	// Untuk preload setelah create
	mockProduct := &model.Product{ID: uuid.New(), NamaProduk: "Produk Tes"}
	mockRepo.On("FindProductByID", mock.Anything).Return(mockProduct, nil)

	// Harapkan publisher dipanggil
	mockPublisher.On("Publish", "product.created", mockProduct).Return(nil)

	// Act
	product, err := service.CreateProduct(inputDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Produk Tes", product.NamaProduk)
	mockRepo.AssertCalled(t, "CreateProduct", mock.Anything)
	mockPublisher.AssertCalled(t, "Publish", "product.created", mockProduct)
}

func TestCreateProduct_CategoryNotFound(t *testing.T) {
	// Arrange
	service, mockRepo, mockPublisher := setupProductServiceTest()

	inputDTO := dto.CreateProductDTO{CategoryID: uuid.New()}

	// "Program" mock-nya (gagal di cek pertama)
	mockRepo.On("FindCategoryByID", inputDTO.CategoryID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	product, err := service.CreateProduct(inputDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "Category ID tidak valid", err.Error())
	// Pastikan tidak ada yg di-publish jika gagal
	mockPublisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
}

func TestDeleteProduct_Success(t *testing.T) {
	// Arrange
	service, mockRepo, mockPublisher := setupProductServiceTest()
	productID := uuid.New()
	mockProduct := &model.Product{ID: productID, NamaProduk: "Akan dihapus"}

	// "Program" mock-nya
	mockRepo.On("FindProductByID", productID).Return(mockProduct, nil)
	mockRepo.On("DeleteProduct", productID).Return(nil)
	mockPublisher.On("Publish", "product.deleted", mockProduct).Return(nil)

	// Act
	err := service.DeleteProduct(productID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "DeleteProduct", productID)
	mockPublisher.AssertCalled(t, "Publish", "product.deleted", mockProduct)
}

func TestUpdateProduct_Success(t *testing.T) {
	// Arrange
	service, mockRepo, mockPublisher := setupProductServiceTest()

	productID := uuid.New()
	inputDTO := dto.UpdateProductDTO{
		NamaProduk: "Produk Update",
	}

	// Buat data palsu "sebelum" update
	mockProduct := &model.Product{
		ID:         productID,
		SKU:        "LAMA-001",
		NamaProduk: "Produk Lama",
	}

	// Buat data palsu "setelah" update
	// (GORM Save akan mengembalikan data yang sudah di-update)
	updatedProduct := &model.Product{
		ID:         productID,
		SKU:        "LAMA-001",
		NamaProduk: "Produk Update", // Nama sudah berubah
	}

	// "Program" mock-nya
	mockRepo.On("FindProductByID", productID).Return(mockProduct, nil).Once() // Panggilan pertama (untuk cek)
	mockRepo.On("UpdateProduct", mock.AnythingOfType("*model.Product")).Return(updatedProduct, nil)

	// Untuk preload setelah update
	mockRepo.On("FindProductByID", productID).Return(updatedProduct, nil).Once() // Panggilan kedua (untuk preload event)

	// Harapkan publisher dipanggil
	mockPublisher.On("Publish", "product.updated", updatedProduct).Return(nil)

	// Act
	product, err := service.UpdateProduct(productID, inputDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Produk Update", product.NamaProduk) // Pastikan data baru yg kembali
	mockRepo.AssertCalled(t, "UpdateProduct", mock.Anything)
	mockPublisher.AssertCalled(t, "Publish", "product.updated", updatedProduct)
}

func TestCreateCategory_Success(t *testing.T) {
	// Arrange
	service, mockRepo, mockPublisher := setupProductServiceTest()

	inputDTO := dto.CreateCategoryDTO{
		NamaKategori: "Kategori Tes",
	}

	// "Program" mock-nya
	mockRepo.On("CreateCategory", mock.AnythingOfType("*model.Category")).Return(nil)
	mockPublisher.On("Publish", "category.created", mock.Anything).Return(nil)

	// Act
	category, err := service.CreateCategory(inputDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Kategori Tes", category.NamaKategori)
	mockRepo.AssertCalled(t, "CreateCategory", mock.Anything)
	mockPublisher.AssertCalled(t, "Publish", "category.created", mock.Anything)
}
