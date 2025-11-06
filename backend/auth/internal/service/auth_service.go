package service

import (
	"errors"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService interface
type AuthService interface {
	Register(input dto.RegisterUserDTO) (*model.User, error)
	Login(input dto.LoginUserDTO) (string, error)
}

// authService struct
type authService struct {
	repo      repository.AuthRepository
	jwtSecret string // Ambil dari config
}

// NewAuthService "constructor"
// Menerima repo dan jwtSecret (dari config)
func NewAuthService(repo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Implementasi Register
func (s *authService) Register(input dto.RegisterUserDTO) (*model.User, error) {
	// 1. Cek apakah role valid
	_, err := s.repo.FindRoleByID(input.RoleID)
	if err != nil {
		return nil, errors.New("Role ID tidak valid")
	}

	// 2. Cek apakah email sudah terdaftar
	_, err = s.repo.FindUserByEmail(input.Email)
	if err == nil { // user ditemukan (sudah ada)
		return nil, errors.New("Email sudah terdaftar")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // Error lain selain "tidak ketemu"
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Buat objek User baru
	newUser := &model.User{
		Nama:         input.Nama,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		RoleID:       input.RoleID,
	}

	// 5. Simpan ke database
	err = s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	// Kosongkan password hash sebelum dikembalikan
	newUser.PasswordHash = ""
	return newUser, nil
}

// Implementasi Login
func (s *authService) Login(input dto.LoginUserDTO) (string, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("Email atau password salah")
	}

	// 2. Bandingkan password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		// Jika error (password tidak cocok)
		return "", errors.New("Email atau password salah")
	}

	// 3. Jika berhasil, generate JWT
	tokenString, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper untuk generate JWT
func (s *authService) generateJWT(user *model.User) (string, error) {
	// Buat claims (data yang disimpan di dalam token)
	claims := jwt.MapClaims{
		"sub":   user.ID, // "Subject" (standard claim)
		"role":  user.Role.NamaPeran,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token berlaku 1 jam
		"iat":   time.Now().Unix(),                    // Issued at
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
