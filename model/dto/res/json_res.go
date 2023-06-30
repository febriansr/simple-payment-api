package res

import "github.com/gin-gonic/gin"

type JsonResponse struct {
	c              *gin.Context
	httpStatusCode int
	response       ApiResponse
}

func (j *JsonResponse) Send() {
	j.c.JSON(j.httpStatusCode, j.response)
}

func (j *JsonResponse) Get() (int, ApiResponse) {
	return j.httpStatusCode, j.response
}

func NewSuccessJsonResponse(c *gin.Context, data any) AppHttpResponse {
	httpStatusCode, res := NewSuccessMessage(data)
	return &JsonResponse{
		c,
		httpStatusCode,
		res,
	}
}

func NewErrorJsonResponse(c *gin.Context, err error) AppHttpResponse {
	httpStatusCode, res := NewFailedMessage(err)
	return &JsonResponse{
		c,
		httpStatusCode,
		res,
	}
}

func NewJsonResponse(c *gin.Context, httpStatusCode int, res ApiResponse) AppHttpResponse {
	return &JsonResponse{
		c,
		httpStatusCode,
		res,
	}
}
