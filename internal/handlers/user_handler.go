package handler

import (
	"errors"
	"net/http"

	service "backend-task/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {

	return &UserHandler{Service: s}
}

type createUserReq struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var bodies []createUserReq
	if err := c.ShouldBindJSON(&bodies); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "invalid request body"})
		return
	}

	for _, body := range bodies {

		_, err := h.Service.CreateUser(body.Name, body.Email, body.DateOfBirth)
		if err != nil {

			respondErr(c, err)
			return
		}
	}

	c.JSON(http.StatusCreated, bodies)
}

type updateUserReq struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	// Group Is Intentionally Omitted To Enforce Read-Only.
}

func (h *UserHandler) UpdateUser(c *gin.Context) {

	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Id !"})
		return
	}

	_, err = h.Service.GetUser(idStr)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "User Not Found !"})
		return
	}

	var body updateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()}) // "Invalid Request Body"
		return
	}

	u, err := h.Service.UpdateUser(idStr, body.Name, body.Email)
	if err != nil {

		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	id := c.Param("id")

	u, err := h.Service.GetUser(id)
	if err != nil {

		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) QueryUsers(c *gin.Context) {

	group := c.Query("group")

	users, err := h.Service.ListUsers(group)
	if err != nil {

		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func respondErr(c *gin.Context, err error) {

	if err == nil {
		return
	}

	// Custom Errors With Error() string.
	type codedError interface {
		Error() string
	}

	if ae, ok := err.(codedError); ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": ae.Error()})
		return
	}

	// Handle GORM Not Found Explicitly.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "record not found"})
		return
	}

	// Default: 500 Internal Server Error.
	c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "internal server error"})
}
