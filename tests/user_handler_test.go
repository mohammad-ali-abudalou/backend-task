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
	mocks "backend-task/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndAutoGroupHandler(t *testing.T) {

	mockService := new(mocks.UserService)

	// Use A Valid Past DOB.
	dateOfBirth := time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC)

	// Mock CreateUser.
	mockService.On("CreateUser", "Abudalou", "abudalou@test.com", "2000-01-04").
		Return(&models.User{
			Name:        "Abudalou",
			Email:       "abudalou@test.com",
			DateOfBirth: dateOfBirth,
		}, nil)

	// Mock ListUsersByFilter.
	mockService.On("ListUsersByFilter", "").Return([]models.User{
		{
			Name:        "Abudalou",
			Email:       "abudalou@test.com",
			DateOfBirth: dateOfBirth,
		},
	}, nil)

	// Setup Router With Mock Service.
	route := router.SetupRoutersWithService(mockService)

	users := []map[string]string{
		{"name": "Abudalou", "email": "abudalou@test.com", "date_of_birth": "2000-01-04"},
	}

	body, err := json.Marshal(users)
	assert.NoError(t, err)

	// POST /users
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	route.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Abudalou")

	// GET /users
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	resp = httptest.NewRecorder()
	route.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Abudalou")

	mockService.AssertExpectations(t)
}
