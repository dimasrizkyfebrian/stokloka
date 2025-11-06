package handler

import (
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler interface
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

// authHandler struct
type authHandler struct {
	authService service.AuthService
	validate    *validator.Validate // Tambahkan validator
}

// NewAuthHandler "constructor"
func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
		validate:    validator.New(), // Inisialisasi validator
	}
}

// Implementasi Register
func (h *authHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterUserDTO

	// 1. Parse JSON body ke struct DTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Input JSON tidak valid",
			"error":   err.Error(),
		})
	}

	// 2. Validasi DTO
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	// 3. Panggil service
	user, err := h.authService.Register(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(), // Pesan error dari service (misal: "Email sudah terdaftar")
		})
	}

	// 3. Kirim response sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data":    user, // user tanpa password hash
	})
}

// Implementasi Login
func (h *authHandler) Login(c *fiber.Ctx) error {
	var input dto.LoginUserDTO

	// 1. Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Input JSON tidak valid",
			"error":   err.Error(),
		})
	}

	// 2. Validasi DTO
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	// 3. Panggil service
	token, err := h.authService.Login(input)
	if err != nil {
		// Error dari service (misal: "Email atau password salah")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// 4. Kirim response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"data":    dto.AuthResponse{Token: token}, // Kirim dalam format DTO
	})
}
