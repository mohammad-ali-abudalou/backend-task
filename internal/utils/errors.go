package utils

import (
	"backend-task/internal/constants"
	models "backend-task/internal/user/models"
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

// ---------------- Predefined Constructors ----------------

func NewBadRequest(err error) error {
	return models.ErrorResponse{Code: constants.StatusBadRequest, Message: err.Error()}
}

func NewNotFound(err error) error {
	return models.ErrorResponse{Code: constants.StatusNotFound, Message: err.Error()}
}

func NewInternalError(err error) error {
	return models.ErrorResponse{Code: constants.StatusInternalServerError, Message: err.Error()}
}

// ---------------- Gin Error Responder ----------------

func RespondError(context *gin.Context, err error) {

	if err == nil {

		return
	}

	var apiErr models.ErrorResponse

	// Use APIError If Possible
	if errors.As(err, &apiErr) {

		context.JSON(apiErr.Code, apiErr)
		return
	}

	switch {

	case errors.Is(err, ErrRecordNotFound):
		context.JSON(constants.StatusNotFound, models.ErrorResponse{Code: constants.StatusNotFound, Message: ErrRecordNotFound.Error()})

	case errors.Is(err, ErrUserNotFound):
		context.JSON(constants.StatusNotFound, models.ErrorResponse{Code: constants.StatusNotFound, Message: ErrUserNotFound.Error()})

	case errors.Is(err, ErrInvalidID):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrInvalidID.Error()})

	case errors.Is(err, ErrNameCannotBeEmpty):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrNameCannotBeEmpty.Error()})

	case errors.Is(err, ErrInvalidEmailFormat):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrInvalidEmailFormat.Error()})

	case errors.Is(err, ErrEmailAlreadyExists):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrEmailAlreadyExists.Error()})

	case errors.Is(err, ErrNameIsRequired):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrNameIsRequired.Error()})

	case errors.Is(err, ErrDateOfBirthFormat):
		context.JSON(constants.StatusBadRequest, models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrDateOfBirthFormat.Error()})

	case errors.Is(err, ErrDateOfBirthCannotBeFuture):
		context.JSON(constants.StatusBadRequest,
			models.ErrorResponse{Code: constants.StatusBadRequest, Message: ErrDateOfBirthCannotBeFuture.Error()})

	default:
		context.JSON(constants.StatusInternalServerError, models.ErrorResponse{Code: constants.StatusInternalServerError, Message: ErrInternalError.Error()})
	}
}
