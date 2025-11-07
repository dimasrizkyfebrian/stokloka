package repository

import (
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// AuthRepositoryMock adalah struct mock untuk AuthRepository
type AuthRepositoryMock struct {
	mock.Mock
}

// Implementasi fungsi mock
func (m *AuthRepositoryMock) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *AuthRepositoryMock) FindUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *AuthRepositoryMock) FindUserByID(userID uuid.UUID) (*model.User, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *AuthRepositoryMock) FindRoleByID(roleID uint) (*model.Role, error) {
	args := m.Called(roleID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}
