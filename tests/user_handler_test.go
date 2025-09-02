package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-task/internal/router"
	models "backend-task/internal/user/models"
	mocks "backend-task/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler(testingT *testing.T) {

	gin.SetMode(gin.TestMode)

	mockService := new(mocks.UserService)

	// Use A Valid Past DOB.
	dateOfBirth := time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC)

	// Create test UUID.
	testUUID := uuid.New()

	name := "Abudalou1"
	email := "abudalou1@test.com"

	// Mock CreateUser.
	mockService.On("CreateUser", "Abudalou", "abudalou@test.com", "2000-01-04").
		Return(&models.User{
			ID:          testUUID,
			Name:        "Abudalou",
			Email:       "abudalou@test.com",
			DateOfBirth: dateOfBirth,
		}, nil)

	// Mock UpdateUser.
	mockService.On("UpdateUser", testUUID.String(), &name, &email).
		Return(&models.User{
			ID:          testUUID,
			Name:        name,
			Email:       email,
			DateOfBirth: dateOfBirth,
		}, nil)

	// Mock GetUserByID.
	mockService.On("GetUserByID", testUUID.String()).
		Return(&models.User{
			ID:          testUUID,
			Name:        "Abudalou1",
			Email:       "abudalou1@test.com",
			DateOfBirth: dateOfBirth,
		}, nil)

	// Mock ListUsersByFilter.
	mockService.On("ListUsersByFilter", "adult-1").
		Return([]*models.User{
			{
				ID:          testUUID,
				Name:        "Abudalou1",
				Email:       "abudalou1@test.com",
				DateOfBirth: dateOfBirth,
			},
		}, nil)

	// Setup Router With Mock Service.
	route := router.SetupRoutersWithService(mockService)

	// 1. CreateUser : POST /users
	users := []map[string]string{
		{"name": "Abudalou", "email": "abudalou@test.com", "date_of_birth": "2000-01-04"},
	}

	creatBody, err := json.Marshal(users)
	assert.NoError(testingT, err)

	creatReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(creatBody))
	creatReq.Header.Set("Content-Type", "application/json")

	creatResp := httptest.NewRecorder()
	route.ServeHTTP(creatResp, creatReq)

	assert.Equal(testingT, http.StatusCreated, creatResp.Code)
	assert.Contains(testingT, creatResp.Body.String(), "Abudalou")

	// Extract user ID from response
	var createdUsers []models.User
	json.Unmarshal(creatResp.Body.Bytes(), &createdUsers)
	createdUser := createdUsers[0]

	// 2. UpdateUser : PATCH /users/:id
	updatePayload := models.UpdateUserReq{
		Name:  &name,
		Email: &email,
	}

	updateBody, err := json.Marshal(updatePayload)
	assert.NoError(testingT, err)

	updateReq := httptest.NewRequest(http.MethodPatch, "/users/"+createdUser.ID.String(), bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp := httptest.NewRecorder()
	route.ServeHTTP(updateResp, updateReq)

	assert.Equal(testingT, http.StatusOK, updateResp.Code)
	assert.Contains(testingT, updateResp.Body.String(), "Abudalou")

	// 3. GetUserByID : GET /users/:id
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUUID.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	route.ServeHTTP(resp, req)

	assert.Equal(testingT, http.StatusOK, resp.Code)
	assert.Contains(testingT, resp.Body.String(), "Abudalou")

	// 4. ListUsersByFilter : GET /users?group=adult-1
	req = httptest.NewRequest(http.MethodGet, "/users?group=adult-1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp = httptest.NewRecorder()
	route.ServeHTTP(resp, req)

	assert.Equal(testingT, http.StatusOK, resp.Code)
	assert.Contains(testingT, resp.Body.String(), "Abudalou1")

	// Verify All Expectations.
	mockService.AssertExpectations(testingT)
}
