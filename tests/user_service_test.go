package tests

import (
	"testing"
	"time"

	"backend-task/internal/user/models"
	mocks "backend-task/tests/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(testingT *testing.T) {
	mockService := new(mocks.UserService)

	user := &models.User{
		ID:          uuid.New(),
		Name:        "Abudalou",
		Email:       "abudalou@test.com",
		DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	// Create User :
	mockService.On("CreateUser", user.Name, user.Email, user.DateOfBirth.String()).Return(user, nil)
	createdUser, err := mockService.CreateUser(user.Name, user.Email, user.DateOfBirth.String())
	assert.NoError(testingT, err)

	// Update User :
	name := "Updated Name"
	email := "updated@test.com"
	mockService.On("UpdateUser", createdUser.ID.String(), &name, &email).Return(user, nil)
	_, err = mockService.UpdateUser(createdUser.ID.String(), &name, &email)
	assert.NoError(testingT, err)

	// Get User By ID :
	mockService.On("GetUserByID", createdUser.ID.String()).Return(user, nil)
	result, err := mockService.GetUserByID(createdUser.ID.String())
	assert.NoError(testingT, err)
	assert.Equal(testingT, user, result)

	// List Users By Filter :
	mockService.On("ListUsersByFilter", "adult-1").Return([]*models.User{user}, nil)
	users, err := mockService.ListUsersByFilter("adult-1")
	assert.NoError(testingT, err)
	assert.Len(testingT, users, 1)
	assert.Equal(testingT, user, users[0])

	// Verify All Expectations :
	mockService.AssertExpectations(testingT)
}

func TestListUsersByFilterService(t *testing.T) {

	mockService := new(mocks.UserService)

	users := []*models.User{
		{ID: uuid.New(), Name: "Abudalou", Email: "abudalou@test.com", DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC)},
		{ID: uuid.New(), Name: "Abudalou1", Email: "abudalou1@test.com", DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC)},
	}

	mockService.On("ListUsersByFilter", "adult-1").Return(users, nil)

	result, err := mockService.ListUsersByFilter("adult-1")
	assert.NoError(t, err)

	assert.Len(t, result, len(users))
	assert.Equal(t, users, result)

	mockService.AssertExpectations(t)
}
