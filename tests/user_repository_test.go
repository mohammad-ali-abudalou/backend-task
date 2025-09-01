package tests

import (
	"context"
	"testing"
	"time"

	"backend-task/internal/user/models"
	mocks "backend-task/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewUserWithMockRepo(t *testing.T) {

	ctx := context.Background()
	mockRepo := new(mocks.UserRepository)

	// Use Realistic DOB In The Past.
	user := &models.User{
		Name:        "Abudalou",
		Email:       "abudalou@test.com",
		DateOfBirth: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
		Group:       "adult-1",
	}

	// Setup Expected Call.
	mockRepo.On("CreateNewUser", ctx, user).Return(nil)

	// Call the method.
	err := mockRepo.CreateNewUser(ctx, user)
	assert.NoError(t, err)

	// Verify That Expectations Were Met.
	mockRepo.AssertExpectations(t)
}
