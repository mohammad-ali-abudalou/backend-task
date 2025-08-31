package tests

import (
	context "context"
	"testing"
	"time"

	"backend-task/internal/user/models"
	repository "backend-task/internal/user/repository/mocks"
)

func TestCreateAndGetUserRepository(t *testing.T) {

	ctx := context.Background()
	mockRepo := new(repository.UserRepository)

	user := &models.User{
		Name:        "Abudalou",
		Email:       "Abudalou@test.com",
		DateOfBirth: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	// Setup the expected call on the mock
	mockRepo.On("CreateNewUser", ctx, user).Return(nil)

	// Call the method
	err := mockRepo.CreateNewUser(ctx, user)
	if err != nil {
		t.Log("Error:", err.Error())
	}

	// Verify that expectations were met
	mockRepo.AssertExpectations(t)
}
