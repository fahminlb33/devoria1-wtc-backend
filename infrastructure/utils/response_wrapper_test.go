package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWrapResponse(t *testing.T) {
	result := utils.WrapResponse(http.StatusOK, "OK", nil)
	assert.Equal(t, result.HttpStatus, http.StatusOK)
}

func TestWriteResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.WriteResponse(c, utils.Response{
		HttpStatus: http.StatusOK,
		Message:    "OK",
		Data:       nil,
	})

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
}

func TestWriteAbortResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.WriteAbortResponse(c, utils.Response{
		HttpStatus: http.StatusBadRequest,
		Message:    "Bad request",
		Data:       nil,
	})

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}
