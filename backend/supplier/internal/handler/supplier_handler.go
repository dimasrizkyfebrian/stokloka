package handler

import (
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SupplierHandler interface {
	Create(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type supplierHandler struct {
	service  service.SupplierService
	validate *validator.Validate
}

func NewSupplierHandler(service service.SupplierService) SupplierHandler {
	return &supplierHandler{service: service, validate: validator.New()}
}

func (h *supplierHandler) Create(c *fiber.Ctx) error {
	var input dto.CreateSupplierDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.service.Create(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"status": "success", "data": res})
}

func (h *supplierHandler) FindAll(c *fiber.Ctx) error {
	res, err := h.service.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func (h *supplierHandler) FindByID(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	res, err := h.service.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Supplier not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func (h *supplierHandler) Update(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var input dto.UpdateSupplierDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.service.Update(id, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func (h *supplierHandler) Delete(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
