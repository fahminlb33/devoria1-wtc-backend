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
	"github.com/stretchr/testify/suite"
)

// --- BasicAuthMiddleware

type BasicAuthMiddlewareSuite struct {
	suite.Suite
}

func (s *BasicAuthMiddlewareSuite) TestBasicAuthWithoutAuthHeader() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	authentication.BasicAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, w.Result().StatusCode)
	assert.Equal(s.T(), "Basic realm=DEVORIA", w.Result().Header.Get("WWW-Authenticate"))
	assert.Equal(s.T(), "Missing authorization header", response["message"])
}

func (s *BasicAuthMiddlewareSuite) TestBasicAuthWithInvalidAuthHeader() {
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
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, w.Result().StatusCode)
	assert.Equal(s.T(), "Basic realm=DEVORIA", w.Result().Header.Get("WWW-Authenticate"))
	assert.Equal(s.T(), "Authorization header is not Basic", response["message"])
}

func (s *BasicAuthMiddlewareSuite) TestBasicAuthWithInvalidAuthValue() {
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
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, w.Result().StatusCode)
	assert.Equal(s.T(), "Basic realm=DEVORIA", w.Result().Header.Get("WWW-Authenticate"))
	assert.Contains(s.T(), response["message"], "illegal base64 data")
}

func (s *BasicAuthMiddlewareSuite) TestBasicAuthWithInvalidCredentials() {
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
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, w.Result().StatusCode)
	assert.Equal(s.T(), "Basic realm=DEVORIA", w.Result().Header.Get("WWW-Authenticate"))
	assert.Contains(s.T(), response["message"], "Invalid credentials")
}

func (s *BasicAuthMiddlewareSuite) TestBasicAuthPositive() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	config.GlobalConfig.Authentication.BasicUsername = "user"
	config.GlobalConfig.Authentication.BasicPassword = "password"

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Basic dXNlcjpwYXNzd29yZA==")
	c.Request = request

	authentication.BasicAuthMiddleware()(c)

	assert.True(s.T(), c.GetBool("BASIC_AUTHENTICATED"))
}

func TestBasicAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(BasicAuthMiddlewareSuite))
}
