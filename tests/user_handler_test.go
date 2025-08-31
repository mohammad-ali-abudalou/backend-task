package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-task/internal/router"
	"backend-task/internal/user/models"
	"backend-task/internal/user/services/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndAutoGroupHandler(t *testing.T) {

	mockService := new(mocks.UserService)

	dateOfBirth := time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC)

	// Mock CreateUser :
	mockService.On("CreateUser", "Abudalou", "Abudalou@test.com", "2025-01-04").
		Return(&models.User{
			Name:        "Abudalou",
			Email:       "Abudalou@test.com",
			DateOfBirth: dateOfBirth,
		}, nil)

	// Mock ListUsersByFilter :
	mockService.On("ListUsersByFilter", "").Return([]models.User{
		{
			Name:        "Abudalou",
			Email:       "Abudalou@test.com",
			DateOfBirth: dateOfBirth,
		},
	}, nil)

	route := router.SetupRoutersWithService(mockService)

	users := []map[string]string{
		{"name": "Abudalou", "email": "Abudalou@test.com", "date_of_birth": "2025-01-04"},
	}

	body, err := json.Marshal(users)
	assert.NoError(t, err)

	// Simulate To POST /user :
	request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	route.ServeHTTP(responseRecorder, request)
	defer request.Body.Close()

	assert.Equal(t, http.StatusCreated, responseRecorder.Code, "Expected 201 Created")
	assert.Contains(t, responseRecorder.Body.String(), "Abudalou")

	// Simulate To GET /users :
	request = httptest.NewRequest(http.MethodGet, "/users", nil)
	responseRecorder = httptest.NewRecorder()
	route.ServeHTTP(responseRecorder, request)
	defer request.Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected 200 OK")
	assert.Contains(t, responseRecorder.Body.String(), "Abudalou")

	mockService.AssertExpectations(t)
}
