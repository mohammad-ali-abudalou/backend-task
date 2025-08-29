package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-task/config"
	"backend-task/internal/models"
	"backend-task/internal/router"

	"gorm.io/gorm"
)

type createReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	DOB   string `json:"date_of_birth"`
}

var gormDB *gorm.DB

func ConnectDB() {

	// Connect DB :
	config.ConnectDatabase()

	// Auto Migrate Schema :
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration Failed : %v", err)
	}

	gormDB = config.DB
}

func TestCreateAndAutoGroup(t *testing.T) {

	// Connect DB :
	ConnectDB()

	// Setup Router :
	r := router.SetupRouter(config.DB)

	// create 4 adults -> adult-1 gets 3, adult-2 gets 1
	bodies := []createReq{
		{"A1", "a1@test.com", "1988-01-01"},
		{"A2", "a2@test.com", "1985-01-01"},
		{"A3", "a3@test.com", "1990-01-01"},
		{"A4", "a4@test.com", "1989-01-01"},
	}

	for _, b := range bodies {

		buf, _ := json.Marshal(map[string]string{"name": b.Name, "email": b.Email, "date_of_birth": b.DOB})
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {

			t.Fatalf("expected 201 got %d", w.Code)
		}
	}

	// query adult-1
	req := httptest.NewRequest(http.MethodGet, "/users?group=adult-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list expected 200 got %d", w.Code)
	}
}
