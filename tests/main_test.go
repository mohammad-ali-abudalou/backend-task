package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend-task/config"
	"backend-task/internal/models"
	"backend-task/internal/router"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

// Setup Test DB & Router :
func setupTestEnv() http.Handler {

	// Load .env.test :
	if err := godotenv.Load(".env.test"); err != nil {
		log.Println("No .env.test File Found !")
	}

	// Connect DB :
	config.ConnectDatabase()

	// Auto Migrate Schema :
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration Failed : %v", err)
	}

	gormDB = config.DB

	// Setup Gin Router :
	return router.SetupRouter(config.DB)
}

func TestMain(m *testing.M) {

	code := m.Run()
	os.Exit(code)
}

func TestCreateAndAutoGroup(t *testing.T) {

	r := setupTestEnv()

	// Users To create — ** Send As A Single JSON Array** :
	users := []map[string]string{
		{"name": "Abudalou", "email": "Abudalou@test.com", "date_of_birth": "2025-01-04"},
	}

	// Marshal Request Body :
	body, err := json.Marshal(users)
	assert.NoError(t, err)

	// Simulate POST /user :
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusCreated, w.Code, "Expected 201 Created")

	// Simulate GET /users :
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	assert.Contains(t, w.Body.String(), "Abudalou")
}
