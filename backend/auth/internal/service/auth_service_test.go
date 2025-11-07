package service

import (
	"testing"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/dto"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/model"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper untuk membuat object service dan mock repo
func setupAuthServiceTest() (AuthService, *repository.AuthRepositoryMock) {
	mockRepo := new(repository.AuthRepositoryMock)
	// JWT untuk test
	jwtSecret := "test_secret"
	authService := NewAuthService(mockRepo, jwtSecret)
	return authService, mockRepo
}

// --- Test Login ---
func TestLogin_Success(t *testing.T) {
	// Persiapan
	authService, mockRepo := setupAuthServiceTest()

	// Password hash yang benar
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Buat data user palsu
	mockUser := &model.User{
		Email:        "test@stokloka.com",
		PasswordHash: string(hashedPassword),
		Role:         model.Role{NamaPeran: "admin"},
	}

	// "Program" mock-nya:
	mockRepo.On("FindUserByEmail", "test@stokloka.com").Return(mockUser, nil)

	// Buat DTO input
	loginDTO := dto.LoginUserDTO{
		Email:    "test@stokloka.com",
		Password: "password123",
	}

	// Eksekusi
	token, err := authService.Login(loginDTO)

	// Verifikasi
	assert.NoError(t, err)         // Pastikan tidak ada error
	assert.NotEmpty(t, token)      // Pastikan token-nya tidak kosong
	mockRepo.AssertExpectations(t) // Pastikan mock-nya dipanggil
}

func TestLogin_UserNotFound(t *testing.T) {
	// Persiapan
	authService, mockRepo := setupAuthServiceTest()

	// "Program" mock-nya:
	mockRepo.On("FindUserByEmail", "salah@stokloka.com").Return(nil, gorm.ErrRecordNotFound)

	// Buat DTO input
	loginDTO := dto.LoginUserDTO{
		Email:    "salah@stokloka.com",
		Password: "password123",
	}

	// Eksekusi
	token, err := authService.Login(loginDTO)

	// Verifikasi
	assert.Error(t, err)                                      // Pastikan ada error
	assert.Empty(t, token)                                    // Pastikan token-nya kosong
	assert.Equal(t, "Email atau password salah", err.Error()) // Cek pesan error
}

func TestLogin_WrongPassword(t *testing.T) {
	// Persiapan
	authService, mockRepo := setupAuthServiceTest()

	// Password hash yang benar
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("passwordBENAR"), bcrypt.DefaultCost)

	// Buat data user palsu yang benar
	mockUser := &model.User{
		Email:        "test@stokloka.com",
		PasswordHash: string(hashedPassword),
		Role:         model.Role{NamaPeran: "admin"},
	}

	// "Program" mock-nya:
	mockRepo.On("FindUserByEmail", "test@stokloka.com").Return(mockUser, nil)

	// Buat DTO input dengan password yang salah
	loginDTO := dto.LoginUserDTO{
		Email:    "test@stokloka.com",
		Password: "passwordSALAH", // Password input berbeda
	}

	// Eksekusi
	token, err := authService.Login(loginDTO)

	// Verifikasi
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "Email atau password salah", err.Error())
}

// --- Tes Register ---

func TestRegister_Success(t *testing.T) {
	// Persiapan
	authService, mockRepo := setupAuthServiceTest()

	// Buat data user palsu untuk registrasi
	registerDTO := dto.RegisterUserDTO{
		Nama:     "User Baru",
		Email:    "baru@stokloka.com",
		Password: "password123",
		RoleID:   3,
	}

	// "Program" mock-nya:
	mockRepo.On("FindRoleByID", uint(3)).Return(&model.Role{ID: 3}, nil)                    // Role ditemukan
	mockRepo.On("FindUserByEmail", "baru@stokloka.com").Return(nil, gorm.ErrRecordNotFound) // Email belum ada
	// Untuk 'CreateUser', hanya perlu pastikan dia kembali 'nil' (sukses)
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	// Eksekusi
	user, err := authService.Register(registerDTO)

	// Verifikasi
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "User Baru", user.Nama)
	assert.Equal(t, "baru@stokloka.com", user.Email)
	assert.Empty(t, user.PasswordHash) // Pastikan password hash dikosongkan
}

func TestRegister_EmailExists(t *testing.T) {
	// Persiapan
	authService, mockRepo := setupAuthServiceTest()

	// Buat data user palsu untuk register
	registerDTO := dto.RegisterUserDTO{
		Nama:     "User Baru",
		Email:    "sudah.ada@stokloka.com",
		Password: "password123",
		RoleID:   3,
	}

	// "Program" mock-nya:
	mockRepo.On("FindRoleByID", uint(3)).Return(&model.Role{ID: 3}, nil)                // Role ditemukan
	mockRepo.On("FindUserByEmail", "sudah.ada@stokloka.com").Return(&model.User{}, nil) // Email DITEMUKAN (nil error)

	// Eksekusi
	user, err := authService.Register(registerDTO)

	// Verifikasi
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "Email sudah terdaftar", err.Error())
}
