package middleware

import (
	"errors"
	"strings"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Tipe untuk data user yang disimpan di context
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
		jwtSecret: cfg.JWTSecret, // Ambil secret dari config service ini
	}
}

// Protected adalah middleware utama untuk memvalidasi JWT
func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": "Missing authorization header",
			})
		}

		// Cek format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": "Invalid authorization header format",
			})
		}
		tokenString := parts[1]

		// Parse dan validasi token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": err.Error(),
			})
		}

		// Simpan data user di context Fiber
		c.Locals("user", claims)

		// Lanjut ke handler/middleware berikutnya
		return c.Next()
	}
}

// RequiredRole adalah middleware untuk cek role
func (m *AuthMiddleware) RequiredRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		return c.Next()
	}
}

// validateToken (helper internal)
func (m *AuthMiddleware) validateToken(tokenString string) (*UserClaims, error) {
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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
