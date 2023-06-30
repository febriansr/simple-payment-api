package res

import (
	"errors"
	"net/http"

	"github.com/febriansr/simple-payment-api/model/app_error"
)

const (
	SuccessCode    = "200"
	SuccessMessage = "Success"

	DefaultErrorCode    = "XX"
	DefaultErrorMessage = "Something went wrong"
)

type AppHttpResponse interface {
	Send()
	Get() (int, ApiResponse)
}

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Status
	Data any `json:"data,omitempty"`
}

func NewSuccessMessage(data any) (httpStatusCode int, apiResponse ApiResponse) {
	status := Status{
		Code:    SuccessCode,
		Message: SuccessMessage,
	}
	httpStatusCode = http.StatusOK
	apiResponse = ApiResponse{
		status, data,
	}
	return
}

func NewFailedMessage(err error) (httpStatusCode int, apiResponse ApiResponse) {
	var status Status
	var userError *app_error.AppError
	if errors.As(err, &userError) {
		status = Status{
			Code:    userError.ErrorCode,
			Message: userError.ErrorMessage,
		}
		httpStatusCode = userError.ErrorType
	} else {
		status = Status{
			Code:    DefaultErrorCode,
			Message: DefaultErrorMessage,
		}
		httpStatusCode = http.StatusInternalServerError
	}
	apiResponse = ApiResponse{
		status, nil,
	}

	return
}
