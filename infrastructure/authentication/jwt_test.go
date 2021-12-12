package authentication_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func GenerateRSAKeyPair() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Extract public component.
	publicKey := privateKey.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	publicKeyPacked, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyPacked,
		},
	)

	return base64.StdEncoding.EncodeToString(publicKeyPem), base64.StdEncoding.EncodeToString(privateKeyPem)
}

// --- InitializeJwtAuth

type InitializeJwtAuthSuite struct {
	suite.Suite
}

func (s *InitializeJwtAuthSuite) TestInitializeJwtAuthPositive() {
	publicKey, privateKey := GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey

	_, err := authentication.ConstructJwtAuth()

	assert.Nil(s.T(), err)
}

func (s *InitializeJwtAuthSuite) TestInitializeJwtAuthInvalidPublicKey() {
	_, privateKey := GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = ""
	config.GlobalConfig.Authentication.PrivateKey = privateKey

	_, err := authentication.ConstructJwtAuth()

	assert.Contains(s.T(), err.Error(), "public key")
}

func (s *InitializeJwtAuthSuite) TestInitializeJwtAuthInvalidPrivateKey() {
	publicKey, _ := GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = ""

	_, err := authentication.ConstructJwtAuth()

	assert.Contains(s.T(), err.Error(), "private key")
}

func TestInitializeJwtAuthSuite(t *testing.T) {
	suite.Run(t, new(InitializeJwtAuthSuite))
}

// --- Sign

func TestSign(t *testing.T) {
	publicKey, privateKey := GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey

	result, _ := authentication.ConstructJwtAuth()
	tokenString, err := result.Sign("1", "fahmi")

	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)
}

// --- JwtAuthMiddleware

type JwtAuthMiddlewareSuite struct {
	suite.Suite
}

func (s *JwtAuthMiddlewareSuite) SetupSuite() {
	publicKey, privateKey := GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *JwtAuthMiddlewareSuite) TestJwtAuthMiddlewareWithoutHeader() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest("GET", "/", nil)
	c.Request = request

	result, _ := authentication.ConstructJwtAuth()
	result.JwtAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusUnauthorized)
	assert.Equal(s.T(), "Missing authorization header", response["message"])
}

func (s *JwtAuthMiddlewareSuite) TestJwtAuthMiddlewareWithInvalidHeader() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Basic asas")
	c.Request = request

	result, _ := authentication.ConstructJwtAuth()
	result.JwtAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, w.Result().StatusCode)
	assert.Equal(s.T(), "Authorization header is not Bearer", response["message"])
}

func (s *JwtAuthMiddlewareSuite) TestJwtAuthMiddlewareWithInvalidToken() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer asas")
	c.Request = request

	result, _ := authentication.ConstructJwtAuth()
	result.JwtAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusUnauthorized)
	assert.Contains(s.T(), response["message"], "invalid")
}

func (s *JwtAuthMiddlewareSuite) TestJwtAuthMiddlewareWithExpiredToken() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result, _ := authentication.ConstructJwtAuth()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(result.PrivateKey)

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	c.Request = request

	result.JwtAuthMiddleware()(c)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), w.Result().StatusCode, http.StatusUnauthorized)
	assert.Contains(s.T(), response["message"], "expired")
}

func (s *JwtAuthMiddlewareSuite) TestJwtAuthMiddlewarePositive() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// create JWT middleware
	result, _ := authentication.ConstructJwtAuth()
	tokenString, _ := result.Sign("1", "fahmi")

	// create request
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	c.Request = request

	// execute the middleware
	result.JwtAuthMiddleware()(c)

	// assertion
	assert.True(s.T(), c.GetBool("JWT_AUTHENTICATED"))
}

func TestJwtAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(JwtAuthMiddlewareSuite))
}

// --- GetJwtUser

type GetJwtUserSuite struct {
	suite.Suite
}

func (s *GetJwtUserSuite) TestGetJwtUserNotAuthenticated() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("JWT_AUTHENTICATED", false)

	result, err := authentication.GetJwtUser(c)

	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "not found")
}

func (s *GetJwtUserSuite) TestGetJwtUserPositive() {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	payload := authentication.JwtPayload{
		Username: "fahmi",
	}

	c.Set("JWT_AUTHENTICATED", true)
	c.Set("JWT_PAYLOAD", &payload)

	result, err := authentication.GetJwtUser(c)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}

func TestGetJwtUserSuite(t *testing.T) {
	suite.Run(t, new(GetJwtUserSuite))
}
