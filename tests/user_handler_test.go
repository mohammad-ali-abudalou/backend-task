package tests

import (

    "backend-task/internal/models"
    "backend-task/internal/repository"
    "backend-task/internal/services"
    "backend-task/internal/handlers"
    "github.com/gin-gonic/gin"
    "net/http/httptest"
    "strings"
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTest() (*handlers.UserHandler, *gorm.DB) {

    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.User{})
    repo := repository.NewUserRepository(db)
    service := services.NewUserService(repo, db)
    handler := handlers.NewUserHandler(service)
    return handler, db
}

func TestCreateUser(t *testing.T) {

    handler, _ := setupTest()
    router := gin.Default()
    router.POST("/users", handler.CreateUser)

    body := `{"name":"John","email":"john@test.com","date_of_birth":"2000-01-01"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != 201 {
        t.Fatalf("Expected 201, got %d", w.Code)
    }
}
