package utils

import "github.com/gin-gonic/gin"

type Response struct {
	HttpStatus int
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func WrapResponse(httpStatus int, message string, data interface{}) Response {
	return Response{
		HttpStatus: httpStatus,
		Message:    message,
		Data:       data,
	}
}

func WriteResponse(c *gin.Context, response Response) {
	c.JSON(response.HttpStatus, response)
}

func WriteAbortResponse(c *gin.Context, response Response) {
	c.AbortWithStatusJSON(response.HttpStatus, response)
}
