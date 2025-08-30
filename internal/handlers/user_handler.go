package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"

	service "backend-task/internal/services"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {

	return &UserHandler{Service: s}
}

type CreateUserReq struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}

// @Summary Create one or more users
// @Description Creates new users and assigns them to groups automatically (up to 3 per group).
// @Tags users
// @Accept json
// @Produce json
// @Param users body []CreateUserReq true "User info array"
// @Success 201 {array} CreateUserReq
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {

	var bodies []CreateUserReq
	if err := c.ShouldBindJSON(&bodies); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Code": http.StatusBadRequest, "Error": " Invalid Request Body :( .."})
		return
	}

	for _, body := range bodies {

		_, err := h.Service.CreateUser(body.Name, body.Email, body.DateOfBirth)
		if err != nil {

			RespondError(c, err)
			return
		}
	}

	c.JSON(http.StatusCreated, bodies)
}

type UpdateUserReq struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	// To Enforce Read-Only, The Group Is Purposefully Left Out.
}

// @Summary Update a user
// @Description Update user name and/or email by ID (group cannot be updated manually).
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserReq true "User info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid request or ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {

	userId := c.Param("id")
	_, err := uuid.Parse(userId)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Code": http.StatusBadRequest, "Error": " Invalid Id :( .. "})
		return
	}

	_, err = h.Service.GetUser(userId)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"Code": http.StatusNotFound, "Error": " User Not Found :( .. "})
		return
	}

	var body UpdateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Code": http.StatusBadRequest, "Error": err.Error()}) // " Invalid Request Body :( .."
		return
	}

	u, err := h.Service.UpdateUser(userId, body.Name, body.Email)
	if err != nil {

		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

// @Summary Get user by ID
// @Description Retrieve a user by their UUID.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {

	id := c.Param("id")

	u, err := h.Service.GetUser(id)
	if err != nil {

		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

// @Summary List users
// @Description List all users, optionally filtered by group (e.g. adult-1, senior-2).
// @Tags users
// @Produce json
// @Param group query string false "Group name"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (h *UserHandler) QueryUsers(c *gin.Context) {

	group := c.Query("group")

	users, err := h.Service.ListUsers(group)
	if err != nil {

		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func RespondError(c *gin.Context, err error) {

	if err == nil {
		return
	}

	type CustomErrors interface {
		Error() string
	}

	if ae, ok := err.(CustomErrors); ok {

		c.JSON(http.StatusBadRequest, gin.H{"Code": http.StatusBadRequest, "Error": ae.Error()})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"Code": http.StatusNotFound, "Error": "Record Not Found :( .. "})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"Code": http.StatusInternalServerError, "Error": "Internal Server Error ): .. "})
}
