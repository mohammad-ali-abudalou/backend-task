package mocks

import (
	"backend-task/internal/models"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

type UserService interface {
	CreateUser(user *models.User) error
	GetUsers() ([]models.User, error)
}

func (m *UserServiceMock) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserServiceMock) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
