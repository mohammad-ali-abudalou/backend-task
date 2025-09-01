package tests

import (
	"testing"
	"time"

	"backend-task/internal/user/models"
	mocks "backend-task/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserService(t *testing.T) {

	mockService := new(mocks.UserService)

	// Use A Realistic DOB.
	user := &models.User{
		Name:        "Abudalou",
		Email:       "abudalou@test.com",
		DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	dateString := user.DateOfBirth.Format("2006-01-02")

	// Setup Mock Correctly With 3 Arguments.
	mockService.On("CreateUser", user.Name, user.Email, dateString).
		Return(user, nil)

	// Call The Mocked Method.
	result, err := mockService.CreateUser(user.Name, user.Email, dateString)

	assert.NoError(t, err)
	assert.Equal(t, user, result)

	// Verify That Expectations Were Met.
	mockService.AssertExpectations(t)
}

func TestListUsersByFilterService(t *testing.T) {

	mockService := new(mocks.UserService)

	users := []models.User{
		{Name: "Abudalou", Email: "abudalou@test.com"},
	}

	mockService.On("ListUsersByFilter", "").Return(users, nil)

	result, err := mockService.ListUsersByFilter("")
	assert.NoError(t, err)
	assert.Equal(t, users, result)

	mockService.AssertExpectations(t)
}
