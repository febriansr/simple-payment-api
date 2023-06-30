package app_error

import (
	"fmt"
	"net/http"
	"strconv"
)

type AppError struct {
	ErrorCode    string
	ErrorMessage string
	ErrorType    int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, status:%s, err: %s", e.ErrorType, e.ErrorCode, e.ErrorMessage)
}

func InvalidError(msg string) error {
	if msg == "" {
		return &AppError{
			ErrorCode:    strconv.Itoa(http.StatusBadRequest),
			ErrorMessage: "invalid input",
			ErrorType:    http.StatusBadRequest,
		}
	} else {
		return &AppError{
			ErrorCode:    strconv.Itoa(http.StatusBadRequest),
			ErrorMessage: msg,
			ErrorType:    http.StatusBadRequest,
		}
	}
}

func DataNotFound(msg string) error {
	if msg == "" {
		return &AppError{
			ErrorMessage: "no data found",
			ErrorCode:    strconv.Itoa(http.StatusNotFound),
			ErrorType:    http.StatusNotFound,
		}
	} else {
		return &AppError{
			ErrorMessage: msg,
			ErrorCode:    strconv.Itoa(http.StatusNotFound),
			ErrorType:    http.StatusNotFound,
		}
	}
}

func Unauthorized(msg string) error {
	if msg == "" {
		return &AppError{
			ErrorMessage: "unauthorized",
			ErrorCode:    strconv.Itoa(http.StatusUnauthorized),
			ErrorType:    http.StatusUnauthorized,
		}
	} else {
		return &AppError{
			ErrorMessage: msg,
			ErrorCode:    strconv.Itoa(http.StatusUnauthorized),
			ErrorType:    http.StatusUnauthorized,
		}
	}
}

func InternalServerError(msg string) error {
	if msg == "" {
		return &AppError{
			ErrorMessage: "internal server error",
			ErrorCode:    strconv.Itoa(http.StatusInternalServerError),
			ErrorType:    http.StatusInternalServerError,
		}
	} else {
		return &AppError{
			ErrorMessage: msg,
			ErrorCode:    strconv.Itoa(http.StatusInternalServerError),
			ErrorType:    http.StatusInternalServerError,
		}
	}
}
