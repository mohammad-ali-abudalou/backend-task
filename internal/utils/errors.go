package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Common API Error Messages :
var (
	ErrNoEnvFileFound                = errors.New("No .env File Found")
	ErrFailedConnectDatabase         = errors.New("Failed To Connect To Database")
	ErrDatabaseConnectedSuccessfully = errors.New("Database Connected Successfully")
	ErrUnsupportedDBDriver           = errors.New("Unsupported DB Driver")
	ErrMigrationFailed               = errors.New("Migration Failed")
	ErrInvalidRequestBody            = errors.New("Invalid Request Body")
	ErrInvalidId                     = errors.New("Invalid Id")
	ErrUserNotFound                  = errors.New("User Not Found")
	ErrNameIsRequired                = errors.New("Name Is Required")
	ErrInvalidEmailFormat            = errors.New("Invalid Email Format")
	ErrDateOfBirthFormat             = errors.New("date_of_birth Must Be YYYY-MM-DD")
	ErrEmailAlreadyExists            = errors.New("Email Already Exists")
	ErrNameCanNotEmpty               = errors.New("Name Cannot Be Empty")
	ErrRecordNotFound                = errors.New("Record Not Found")
	ErrInternalError                 = errors.New("Internal Server Error")
	ErrDateOfBirthCanNotInFuture     = errors.New("date_of_birth Cannot Be In The Future")
)

type APIError struct {
	Code    int
	Message string
}

func (apiError APIError) Error() string {
	return apiError.Message
}

// Predefined Constructors For API Errors.
func NewBadRequest(message string) error {

	return APIError{Code: StatusBadRequest, Message: message}
}

func NewNotFound(message string) error {

	return APIError{Code: StatusNotFound, Message: message}
}

// React Depending On The Type Of Issue, Issue Provides The Relevant HTTP Response.
func RespondError(context *gin.Context, err error) {

	if err == nil {

		return
	}

	// Handle Issue In Custom APIs :
	type CustomErrors interface {
		Error() string
	}

	if customErrors, ok := err.(CustomErrors); ok {

		context.JSON(http.StatusBadRequest, gin.H{"Code": StatusBadRequest, "Error": customErrors.Error()})
		return
	}

	// Handle The "GORM Record Not Found" Issue :
	if errors.Is(err, gorm.ErrRecordNotFound) {

		context.JSON(http.StatusNotFound, gin.H{"Code": StatusNotFound, "Error": ErrRecordNotFound})
		return
	}

	// " Internal Server Error " is the default.
	context.JSON(http.StatusInternalServerError, gin.H{"Code": StatusInternalServerError, "Error": ErrInternalError})
}
