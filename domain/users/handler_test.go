package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
	"github.com/fahminlb33/devoria1-wtc-backend/mocks"
)

// --- Login

type LoginSuite struct {
	suite.Suite
}

func (s *LoginSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *LoginSuite) TestLoginWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	// call handler
	handler.Login(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *LoginSuite) TestLoginPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	payload := users.LoginModel{
		Email:    "foo@bar.com",
		Password: "foobar",
	}
	bodyJson, _ := json.Marshal(&payload)
	bodyReader := strings.NewReader(string(bodyJson))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bodyReader)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Login", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// call handler
	handler.Login(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestLoginSute(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

// --- Register

type RegisterSuite struct {
	suite.Suite
}

func (s *RegisterSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *RegisterSuite) TestRegisterWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	// call handler
	handler.Register(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *RegisterSuite) TestRegisterPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	payload := users.RegisterModel{
		Email:     "foo@bar.com",
		Password:  "foobar",
		FirstName: "foo",
		LastName:  "bar",
	}
	bodyJson, _ := json.Marshal(&payload)
	bodyReader := strings.NewReader(string(bodyJson))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bodyReader)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Register", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// call handler
	handler.Register(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestRegisterSute(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

// --- Profile

type ProfileSuite struct {
	suite.Suite
}

func (s *ProfileSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *ProfileSuite) TestRegisterWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	// call handler
	handler.GetProfile(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusInternalServerError)
	assert.Equal(s.T(), response["message"], "Can't get user from token")
}

func (s *ProfileSuite) TestRegisterPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	payload := users.RegisterModel{
		Email:     "foo@bar.com",
		Password:  "foobar",
		FirstName: "foo",
		LastName:  "bar",
	}
	bodyJson, _ := json.Marshal(&payload)
	bodyReader := strings.NewReader(string(bodyJson))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.UserUseCase)
	handler := users.ConstructUserHandler(gin.New(), useCase, jwtAuth)

	useCase.On("GetProfile", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set token to request
	c.Request = httptest.NewRequest("GET", "/", bodyReader)
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.GetProfile(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestProfileSuite(t *testing.T) {
	suite.Run(t, new(ProfileSuite))
}
