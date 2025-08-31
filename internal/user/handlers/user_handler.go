package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	models "backend-task/internal/user/models"
	service "backend-task/internal/user/services"
	"backend-task/internal/utils"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {

	return &UserHandler{Service: s}
}

// @Summary Create one or more users
// @Description Creates new users and assigns them to groups automatically (up to 3 per group).
// @Tags users
// @Accept json
// @Produce json
// @Param users body []models.CreateUserReq true "User info array"
// @Success 201 {array} models.CreateUserReq
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [post]
func (userHandler *UserHandler) CreateUser(context *gin.Context) {

	var bodies []models.CreateUserReq
	if err := context.ShouldBindJSON(&bodies); err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"Code": utils.StatusBadRequest, "Error": utils.ErrInvalidRequestBody})
		return
	}

	for _, body := range bodies {

		_, err := userHandler.Service.CreateUser(body.Name, body.Email, body.DateOfBirth)
		if err != nil {

			utils.RespondError(context, err)
			return
		}
	}

	context.JSON(http.StatusCreated, bodies)
}

// @Summary Update a user
// @Description Update user name and/or email by ID (group cannot be updated manually).
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserReq true "User info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid request or ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [patch]
func (userHandler *UserHandler) UpdateUser(context *gin.Context) {

	userId := context.Param("id")
	_, err := uuid.Parse(userId)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"Code": utils.StatusBadRequest, "Error": utils.ErrInvalidId})
		return
	}

	_, err = userHandler.Service.GetUserById(userId)
	if err != nil {

		context.JSON(http.StatusNotFound, gin.H{"Code": utils.StatusNotFound, "Error": utils.ErrUserNotFound})
		return
	}

	var body models.UpdateUserReq
	if err := context.ShouldBindJSON(&body); err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"Code": http.StatusBadRequest, "Error": err.Error()}) // " Invalid Request Body "
		return
	}

	u, err := userHandler.Service.UpdateUser(userId, body.Name, body.Email)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(http.StatusOK, u)
}

// @Summary Get user by ID
// @Description Retrieve a user by their UUID.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func (userHandler *UserHandler) GetUserByID(context *gin.Context) {

	userId := context.Param("id")

	u, err := userHandler.Service.GetUserById(userId)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(http.StatusOK, u)
}

// @Summary List users
// @Description List all users, optionally filtered by group (e.g. adult-1, senior-2).
// @Tags users
// @Produce json
// @Param group query string false "Group name"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (userHandler *UserHandler) QueryUsers(context *gin.Context) {

	group := context.Query("group")

	users, err := userHandler.Service.ListUsersByFilter(group)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(http.StatusOK, users)
}
