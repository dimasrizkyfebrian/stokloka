package middleware

import (
	"errors"
	"strings"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Tipe untuk data user yang akan disimpan di context
type UserClaims struct {
	ID    string `json:"sub"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

// AuthMiddleware struct untuk menampung secret
type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware "constructor"
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: cfg.JWTSecret,
	}
}

// Protected adalah middleware utama untuk memvalidasi JWT
func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": "Missing authorization header",
			})
		}

		// 2. Cek format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": "Invalid authorization header format",
			})
		}
		tokenString := parts[1]

		// 3. Parse dan validasi token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": err.Error(),
			})
		}

		// 4. Simpan data user di context Fiber (c.Locals)
		c.Locals("user", claims)

		// 5. Lanjut ke handler/middleware berikutnya
		return c.Next()
	}
}

// RequiredRole adalah middleware untuk cek role (ADMIN, MANAGER, dll)
// Middleware ini HARUS dijalankan SETELAH Protected()
func (m *AuthMiddleware) RequiredRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil data user dari context (yang sudah di-set oleh Protected())
		user, ok := c.Locals("user").(*UserClaims)
		if !ok || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error", "message": "Akses ditolak (data user tidak ditemukan)",
			})
		}

		// Cek role
		if user.Role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error", "message": "Akses ditolak (role tidak memadai)",
			})
		}

		// Role cocok, lanjutkan
		return c.Next()
	}
}

// validateToken (helper internal)
func (m *AuthMiddleware) validateToken(tokenString string) (*UserClaims, error) {
	// Definisikan struct claims yang kita harapkan
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("Token tidak valid atau kedaluwarsa")
	}

	if !token.Valid {
		return nil, errors.New("Token tidak valid")
	}

	// Konversi MapClaims ke UserClaims struct
	userClaims := &UserClaims{
		ID:    (*claims)["sub"].(string),
		Role:  (*claims)["role"].(string),
		Email: (*claims)["email"].(string),
	}

	return userClaims, nil
}
