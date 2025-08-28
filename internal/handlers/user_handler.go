package handlers

import (

    "backend-task/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/google/uuid"
	
	"fmt"

)

type UserHandler struct {

    Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {

    return &UserHandler{Service: s}
}

type UserRequest struct {
    Name        string `json:"name"`
    Email       string `json:"email"`
    DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

type UserUpdate struct {
    Name        string `json:"name"`
    Email       string `json:"email"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var bodies []UserRequest
	if err := c.ShouldBindJSON(&bodies); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, body := range bodies {
	
		fmt.Println("Name:", body.Name);
		fmt.Println("Email:", body.Email);
		fmt.Println("DateOfBirth:", body.DateOfBirth);
		
		_, err := h.Service.CreateUser(body.Name, body.Email, body.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, bodies)
}


func (h *UserHandler) GetUser(c *gin.Context) {

    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
	
    if err != nil {
	
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    user, err := h.Service.Repo.GetByID(id)
    if err != nil {
	
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {

    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    user, err := h.Service.Repo.GetByID(id)
    if err != nil {
	
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    var body UserUpdate;
	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
    if err := h.Service.UpdateUser(user, body.Name, body.Email); err != nil {
	
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsersByGroup(c *gin.Context) {

    group := c.Query("group")
    users, err := h.Service.Repo.GetByGroup(group)
    if err != nil {
	
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {

    users, err := h.Service.Repo.GetAll()
    if err != nil {
	
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	
    c.JSON(http.StatusOK, users)
}
