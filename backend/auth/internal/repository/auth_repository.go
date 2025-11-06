package repository

import (
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthRepository interface (kontrak)
type AuthRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByID(userID uuid.UUID) (*model.User, error)
	FindRoleByID(roleID uint) (*model.Role, error)
}

// authRepository struct (implementasi)
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository "constructor"
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Implementasi fungsi
func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	// Preload("Role") untuk mengambil data relasi role-nya
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return &user, err // Akan return gorm.ErrRecordNotFound jika tidak ada
}

func (r *authRepository) FindUserByID(userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *authRepository) FindRoleByID(roleID uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, roleID).Error // Cara singkat find by primary key
	return &role, err
}
