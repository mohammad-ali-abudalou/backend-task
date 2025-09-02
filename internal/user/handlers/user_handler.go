package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	constants "backend-task/internal/constants"
	models "backend-task/internal/user/models"
	UserServiceInterface "backend-task/internal/user/services/interface"
	"backend-task/internal/utils"
)

type UserHandler struct {
	Service UserServiceInterface.UserService
}

func NewUserHandler(s UserServiceInterface.UserService) *UserHandler {

	return &UserHandler{Service: s}
}

// CreateUser godoc
// @Summary Create one or more users.
// @Description Creates new users and assigns them to groups assigned automatically ( up to 3 per group ).
// @Tags users
// @Accept json
// @Produce json
// @Param users body []models.CreateUserReq true "User info array"
// @Success 201 {array} models.User
// @Failure 400 {object} models.ErrorResponse "Invalid request. Possible reasons: email already exists, invalid email format, name is required, date_of_birth must be yyyy-mm-dd, or date_of_birth cannot be in the future."
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users [post]
func (userHandler *UserHandler) CreateUser(context *gin.Context) {

	var bodies []models.CreateUserReq
	if err := context.ShouldBindJSON(&bodies); err != nil {

		utils.RespondError(context, utils.ErrInvalidRequestBody)
		return
	}

	var created []models.User
	for _, body := range bodies {

		user, err := userHandler.Service.CreateUser(body.Name, body.Email, body.DateOfBirth)
		if err != nil {

			utils.RespondError(context, err)
			return
		}

		created = append(created, *user)
	}

	context.JSON(constants.StatusCreated, created)
}

// UpdateUser godoc
// @Summary Update a user.
// @Description Update user information (email and name; group cannot be updated manually).
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserReq true "User info"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse "Invalid request. Possible reasons: invalid ID, email already exists, name cannot be empty, or invalid email format."
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (userHandler *UserHandler) UpdateUser(context *gin.Context) {

	userId := context.Param("id")
	if _, err := uuid.Parse(userId); err != nil {

		utils.RespondError(context, utils.ErrInvalidID)
		return
	}

	// Ensure User Exists :
	if _, err := userHandler.Service.GetUserByID(userId); err != nil {

		utils.RespondError(context, utils.ErrUserNotFound)
		return
	}

	var body models.UpdateUserReq
	if err := context.ShouldBindJSON(&body); err != nil {

		utils.RespondError(context, utils.ErrInvalidRequestBody)
		return
	}

	user, err := userHandler.Service.UpdateUser(userId, body.Name, body.Email)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(constants.StatusOK, user)
}

// GetUserByID godoc
// @Summary Get user by ID.
// @Description Retrieve a user by their unique UUID.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse "Invalid request or invalid ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (userHandler *UserHandler) GetUserByID(context *gin.Context) {

	userId := context.Param("id")
	if _, err := uuid.Parse(userId); err != nil {

		utils.RespondError(context, utils.ErrInvalidID)
		return
	}

	user, err := userHandler.Service.GetUserByID(userId)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(constants.StatusOK, user)
}

// QueryUsers godoc
// @Summary Search users by group / List all users
// @Description Returns a list of users, optionally filtered by group using query parameter (e.g., adult-1, senior-2).
// @Tags users
// @Accept json
// @Produce json
// @Param group query string false "Group name"
// @Success 200 {array} models.User "List of users"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users [get]
func (userHandler *UserHandler) QueryUsers(context *gin.Context) {

	group := context.Query("group")

	users, err := userHandler.Service.ListUsersByFilter(group)
	if err != nil {

		utils.RespondError(context, err)
		return
	}

	context.JSON(constants.StatusOK, users)
}
