package tests

import (
	"testing"
	"time"

	"backend-task/internal/user/models"
	mocks "backend-task/internal/user/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserService(t *testing.T) {

	mockService := new(mocks.UserService)

	user := &models.User{
		Name:        "Abudalou",
		Email:       "Abudalou@test.com",
		DateOfBirth: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	dateString := user.DateOfBirth.Format("2006-01-02")

	// Setup mock
	mockService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&models.User{Name: "Abudalou", Email: "Abudalou@test.com", DateOfBirth: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC)}, nil)

	// Call the method correctly
	_, err := mockService.CreateUser(user.Name, user.Email, dateString)

	if err != nil {
		t.Log("Error:", err.Error())
	}

	// Assert that the mock was called as expected
	mockService.AssertExpectations(t)
}

func TestListUsersByFilterService(t *testing.T) {

	mockService := new(mocks.UserService)

	users := []models.User{
		{Name: "Abudalou", Email: "Abudalou@test.com"},
	}

	mockService.On("ListUsersByFilter", "").Return(users, nil)

	result, err := mockService.ListUsersByFilter("")
	assert.NoError(t, err)
	assert.Equal(t, users, result)

	mockService.AssertExpectations(t)
}
