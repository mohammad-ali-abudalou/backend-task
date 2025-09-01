package utils

import (
	"backend-task/internal/constants"
	"errors"

	"github.com/gin-gonic/gin"
)

// ---------------- Predefined Error Messages ----------------
var (
	ErrNoEnvFileFound                     = errors.New("no .env file found, relying on system environment variables")
	ErrFailedConnectDatabase              = errors.New("failed to connect to database")
	ErrDatabaseConnected                  = errors.New("database connected successfully")
	ErrUnsupportedDBDriver                = errors.New("unsupported db driver")
	ErrMigrationFailed                    = errors.New("migration failed")
	ErrDatabaseSchemaMigratedSuccessfully = errors.New("database schema migrated successfully")
	ErrInvalidRequestBody                 = errors.New("invalid request body")
	ErrInvalidID                          = errors.New("invalid id")
	ErrUserNotFound                       = errors.New("user not found")
	ErrNameIsRequired                     = errors.New("name is required")
	ErrInvalidEmailFormat                 = errors.New("invalid email format")
	ErrDateOfBirthFormat                  = errors.New("date_of_birth must be yyyy-mm-dd")
	ErrEmailAlreadyExists                 = errors.New("email already exists")
	ErrNameCannotBeEmpty                  = errors.New("name cannot be empty")
	ErrRecordNotFound                     = errors.New("record not found")
	ErrInternalError                      = errors.New("internal server error")
	ErrDateOfBirthCannotBeFuture          = errors.New("date_of_birth cannot be in the future")
	ErrFailedToFindGroup                  = errors.New("failed to find group")
	ErrFailedToGetMaxGroupIdx             = errors.New("failed to get max group index")
	ErrFailedToCreateNewGroup             = errors.New("failed to create new group")
)

// ---------------- Error Response ----------------

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// ---------------- Predefined Constructors ----------------

func NewBadRequest(msg string) error {
	return ErrorResponse{Code: constants.StatusBadRequest, Message: msg}
}

func NewNotFound(msg string) error {
	return ErrorResponse{Code: constants.StatusNotFound, Message: msg}
}

func NewInternalError(msg string) error {
	return ErrorResponse{Code: constants.StatusInternalServerError, Message: msg}
}

// ---------------- Gin Error Responder ----------------

func RespondError(context *gin.Context, err error) {

	if err == nil {

		return
	}

	var apiErr ErrorResponse

	// Use APIError If Possible
	if errors.As(err, &apiErr) {

		context.JSON(apiErr.Code, apiErr)
		return
	}

	// Handle GORM Record Not Found
	if errors.Is(err, ErrRecordNotFound) {

		context.JSON(constants.StatusNotFound, ErrorResponse{Code: constants.StatusNotFound, Message: ErrRecordNotFound.Error()})
		return
	}

	// Handle GORM User Not Found
	if errors.Is(err, ErrUserNotFound) {

		context.JSON(constants.StatusNotFound, ErrorResponse{Code: constants.StatusNotFound, Message: ErrUserNotFound.Error()})
		return
	}

	// Handle GORM Invalid ID
	if errors.Is(err, ErrInvalidID) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusNotFound, Message: ErrInvalidID.Error()})
		return
	}

	// Default Internal Server Error
	context.JSON(constants.StatusInternalServerError, ErrorResponse{Code: constants.StatusInternalServerError, Message: ErrInternalError.Error()})
}
