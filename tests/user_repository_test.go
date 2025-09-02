package tests

import (
	"context"
	"testing"
	"time"

	"backend-task/internal/user/models"
	mocks "backend-task/tests/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(testingT *testing.T) {

	currentContext := context.Background()
	mockRepo := new(mocks.UserRepository)

	// Create test UUID.
	testUUID := uuid.New()

	user := &models.User{
		ID:          testUUID,
		Name:        "Abudalou",
		Email:       "abudalou@test.com",
		DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
		Group:       "adult-1",
	}

	// Create New User :
	mockRepo.On("CreateNewUser", currentContext, user).Return(nil)
	err := mockRepo.CreateNewUser(currentContext, user)
	assert.NoError(testingT, err)

	// Update User :
	mockRepo.On("UpdateUser", currentContext, user).Return(nil)
	err = mockRepo.UpdateUser(currentContext, user)
	assert.NoError(testingT, err)

	// Get User By ID :
	mockRepo.On("GetUserByID", currentContext, user.ID).Return(user, nil)
	result, err := mockRepo.GetUserByID(currentContext, user.ID)
	assert.NoError(testingT, err)
	assert.Equal(testingT, user, result)

	// List Users By Group :
	mockRepo.On("ListUsers", currentContext, "adult-1").Return([]*models.User{user}, nil)
	users, err := mockRepo.ListUsers(currentContext, "adult-1")
	assert.NoError(testingT, err)
	assert.Len(testingT, users, 1)
	assert.Equal(testingT, user, users[0])

	// Verify All Expectations :
	mockRepo.AssertExpectations(testingT)
}
