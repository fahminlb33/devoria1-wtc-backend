package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/articles"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
	"github.com/fahminlb33/devoria1-wtc-backend/mocks"
)

// --- FindAll

type FindAllSuite struct {
	suite.Suite
}

func (s *FindAllSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *FindAllSuite) TestFindAllWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	request := httptest.NewRequest("GET", "/", nil)
	queries := request.URL.Query()
	queries.Add("page", "-1")

	request.URL.RawQuery = queries.Encode()
	c.Request = request
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.FindAll(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *FindAllSuite) TestFindAllPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	useCase.On("FindAll", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.FindAll(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestFindAllSuite(t *testing.T) {
	suite.Run(t, new(FindAllSuite))
}

// --- Get

type GetSuite struct {
	suite.Suite
}

func (s *GetSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *GetSuite) TestGetWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	request := httptest.NewRequest("GET", "/", nil)
	queries := request.URL.Query()
	queries.Add("page", "-1")

	request.URL.RawQuery = queries.Encode()
	c.Request = request
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Get(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *GetSuite) TestGetPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Get", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("GET", "/api/v1/articles/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Get(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestGetSuite(t *testing.T) {
	suite.Run(t, new(GetSuite))
}

// --- Get

type DeleteSuite struct {
	suite.Suite
}

func (s *DeleteSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *DeleteSuite) TestDeleteWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	request := httptest.NewRequest("GET", "/", nil)
	queries := request.URL.Query()
	queries.Add("page", "-1")

	request.URL.RawQuery = queries.Encode()
	c.Request = request
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Delete(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *DeleteSuite) TestDeletePositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Delete", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("GET", "/api/v1/articles/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Delete(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestDeleteSuite(t *testing.T) {
	suite.Run(t, new(DeleteSuite))
}

// --- Login

type CreateSuite struct {
	suite.Suite
}

func (s *CreateSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey

	// register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isarticlepublishstatus", articles.IsArticlePublishStatus)
	}
}

func (s *CreateSuite) TestCreateWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	// set token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("POST", "/", nil)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Create(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *CreateSuite) TestCreatePositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Create", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	payload := articles.CreateModel{
		Title:   "tiasdsatle",
		Content: "coasdsadntent",
		Slug:    "slasdsadug",
		Status:  articles.DRAFT,
	}
	bodyJson, _ := json.Marshal(&payload)
	bodyReader := strings.NewReader(string(bodyJson))

	c.Request = httptest.NewRequest("POST", "/", bodyReader)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Create(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestCreateSuite(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

// --- Login

type SaveSuite struct {
	suite.Suite
}

func (s *SaveSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey

	// register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isarticlepublishstatus", articles.IsArticlePublishStatus)
	}
}

func (s *SaveSuite) TestSaveWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	// set token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("POST", "/", nil)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Save(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(s.T(), response["message"], "Validation failed")
}

func (s *SaveSuite) TestSavePositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := new(mocks.ArticleUseCase)
	handler := articles.ConstructArticlesHandler(gin.New(), useCase, jwtAuth)

	useCase.On("Save", c, mock.Anything).Return(utils.WrapResponse(http.StatusOK, "Success", "data here"))

	// set body and token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	payload := articles.SaveModel{
		ArticleId: 1,
		Title:     "tiasdsatle",
		Status:    articles.PUBLISHED,
	}
	bodyJson, _ := json.Marshal(&payload)
	bodyReader := strings.NewReader(string(bodyJson))

	c.Request = httptest.NewRequest("POST", "/", bodyReader)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	handler.Save(c)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusOK)
	assert.Equal(s.T(), response["message"], "Success")
}

func TestSaveSuite(t *testing.T) {
	suite.Run(t, new(SaveSuite))
}

// --- InjectUserId

type InjectUserIdSuite struct {
	suite.Suite
}

func (s *InjectUserIdSuite) TestInjectUserIdWithValidationError() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// call handler
	var userId int = 0
	articles.InjectUserId(c, &userId)

	// assert
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusInternalServerError)
	assert.Equal(s.T(), response["message"], "Can't get user from token")
}

func (s *InjectUserIdSuite) TestInjectUserIdPositive() {
	// initialize gin
	gin.SetMode(gin.TestMode)

	// create new request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// set token to request
	jwtPayload := authentication.JwtPayload{}
	jwtPayload.Id = "1"

	c.Request = httptest.NewRequest("POST", "/", nil)
	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &jwtPayload)

	// call handler
	var userId int = 0
	articles.InjectUserId(c, &userId)

	// assert
	assert.Equal(s.T(), userId, 1)
}

func TestInjectUserIdSuite(t *testing.T) {
	suite.Run(t, new(InjectUserIdSuite))
}
