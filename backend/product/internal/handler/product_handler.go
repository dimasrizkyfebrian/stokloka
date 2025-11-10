package handler

import (
	"github.com/dimasrizkyfebrian/stokloka/product/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Definisikan interface
type ProductHandler interface {
	// Category
	CreateCategory(c *fiber.Ctx) error
	GetAllCategories(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error

	// Unit
	CreateUnit(c *fiber.Ctx) error
	GetAllUnits(c *fiber.Ctx) error
	UpdateUnit(c *fiber.Ctx) error
	DeleteUnit(c *fiber.Ctx) error

	// Product
	CreateProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productHandler struct {
	productService service.ProductService
	validate       *validator.Validate
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{
		productService: productService,
		validate:       validator.New(),
	}
}

// --- Implementasi Category ---
func (h *productHandler) CreateCategory(c *fiber.Ctx) error {
	var input dto.CreateCategoryDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	category, err := h.productService.CreateCategory(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": category})
}

func (h *productHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.productService.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": categories})
}

func (h *productHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	var input dto.UpdateCategoryDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	category, err := h.productService.UpdateCategory(id, input)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Category not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": category})
}

func (h *productHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	if err := h.productService.DeleteCategory(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Category not found"})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

// --- Implementasi Unit ---
func (h *productHandler) CreateUnit(c *fiber.Ctx) error {
	var input dto.CreateUnitDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	unit, err := h.productService.CreateUnit(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": unit})
}

func (h *productHandler) GetAllUnits(c *fiber.Ctx) error {
	units, err := h.productService.GetAllUnits()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": units})
}

func (h *productHandler) UpdateUnit(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	var input dto.UpdateUnitDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	unit, err := h.productService.UpdateUnit(id, input)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Unit not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": unit})
}

func (h *productHandler) DeleteUnit(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	if err := h.productService.DeleteUnit(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Unit not found"})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

// --- Implementasi Product ---
func (h *productHandler) CreateProduct(c *fiber.Ctx) error {
	var input dto.CreateProductDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	product, err := h.productService.CreateProduct(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": product})
}

func (h *productHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": products})
}

func (h *productHandler) GetProductByID(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam) // Konversi string ke UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": product})
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	var input dto.UpdateProductDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	product, err := h.productService.UpdateProduct(id, input)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": product})
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Product not found"})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}
