package authentication_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthWithoutAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	authentication.BasicAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	assert.Equal(t, w.Result().Header.Get("WWW-Authenticate"), "Basic realm=DEVORIA")
	assert.Equal(t, "Missing authorization header", response["message"])
}

func TestBasicAuthWithInvalidAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "invalid")
	c.Request = request

	authentication.BasicAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	assert.Equal(t, w.Result().Header.Get("WWW-Authenticate"), "Basic realm=DEVORIA")
	assert.Equal(t, "Authorization header is not Basic", response["message"])
}

func TestBasicAuthWithInvalidAuthValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Basic hhe")
	c.Request = request

	authentication.BasicAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	assert.Equal(t, w.Result().Header.Get("WWW-Authenticate"), "Basic realm=DEVORIA")
	assert.Contains(t, response["message"], "illegal base64 data")
}

func TestBasicAuthWithInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	config.GlobalConfig.Authentication.BasicUsername = "user"
	config.GlobalConfig.Authentication.BasicPassword = "password"

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Basic dXNlcjI6cGFzc3dvcmQ=")
	c.Request = request

	authentication.BasicAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	assert.Equal(t, w.Result().Header.Get("WWW-Authenticate"), "Basic realm=DEVORIA")
	assert.Contains(t, "Invalid credentials", response["message"])
}

func TestBasicAuthPositive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	config.GlobalConfig.Authentication.BasicUsername = "user"
	config.GlobalConfig.Authentication.BasicPassword = "password"

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Basic dXNlcjpwYXNzd29yZA==")
	c.Request = request

	authentication.BasicAuthMiddleware()(c)

	assert.True(t, c.GetBool("BASIC_AUTHENTICATED"))
}
