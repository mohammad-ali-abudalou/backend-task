package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
