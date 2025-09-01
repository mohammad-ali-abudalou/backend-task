package utils

import (
	"backend-task/internal/constants"
	"errors"
	"fmt"

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

func NewBadRequest(err error) error {
	return ErrorResponse{Code: constants.StatusBadRequest, Message: err.Error()}
}

func NewNotFound(err error) error {
	return ErrorResponse{Code: constants.StatusNotFound, Message: err.Error()}
}

func NewInternalError(err error) error {
	return ErrorResponse{Code: constants.StatusInternalServerError, Message: err.Error()}
}

// ---------------- Gin Error Responder ----------------

func RespondError(context *gin.Context, err error) {

	fmt.Println("RespondError")
	fmt.Println(err)

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

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrInvalidID.Error()})
		return
	}

	// Handle GORM Name Can not Be Empty
	if errors.Is(err, ErrNameCannotBeEmpty) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrNameCannotBeEmpty.Error()})
		return
	}

	// Handle GORM Invalid Email Format
	if errors.Is(err, ErrInvalidEmailFormat) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrInvalidEmailFormat.Error()})
		return
	}

	// Handle GORM Email Already Exists
	if errors.Is(err, ErrEmailAlreadyExists) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrEmailAlreadyExists.Error()})
		return
	}

	// Handle GORM Name Is Required
	if errors.Is(err, ErrNameIsRequired) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrNameIsRequired.Error()})
		return
	}

	// Handle GORM Date Of Birth Format
	if errors.Is(err, ErrDateOfBirthFormat) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrDateOfBirthFormat.Error()})
		return
	}

	fmt.Println(err)
	fmt.Println(ErrDateOfBirthCannotBeFuture)

	// Handle GORM Date Of Birth Can Not Be Future
	if errors.Is(err, ErrDateOfBirthCannotBeFuture) {

		context.JSON(constants.StatusBadRequest, ErrorResponse{Code: constants.StatusBadRequest, Message: ErrDateOfBirthCannotBeFuture.Error()})
		return
	}

	fmt.Println("Default Internal Server ErrorError")
	fmt.Println(err)

	// Default Internal Server Error
	context.JSON(constants.StatusInternalServerError, ErrorResponse{Code: constants.StatusInternalServerError, Message: ErrInternalError.Error()})
}
