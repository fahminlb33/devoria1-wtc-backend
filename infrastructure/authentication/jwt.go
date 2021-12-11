package authentication

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.elastic.co/apm"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

var (
	_privateKey *rsa.PrivateKey
	_publicKey  *rsa.PublicKey
)

type JwtPayload struct {
	jwt.StandardClaims
	Username string `json:"username,omitempty"`
}

// Load public-private key pair from file.
func InitializeJwtAuth() {
	privateKeyContents, privateKeyContentsError := ioutil.ReadFile(config.GlobalConfig.Authentication.PrivateKeyPath)
	if privateKeyContentsError != nil {
		log.Fatal("Error reading private key file: ", privateKeyContentsError)
	}

	decodedPrivateKey, decodedPrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(privateKeyContents)
	if decodedPrivateKeyError != nil {
		log.Fatal("Error parsing private key: ", decodedPrivateKeyError)
	}

	publicKeyContents, publicKeyContentsError := ioutil.ReadFile(config.GlobalConfig.Authentication.PublicKeyPath)
	if publicKeyContentsError != nil {
		log.Fatal("Error reading public key file: ", publicKeyContentsError)
	}

	decodedPublicKey, decodedPublicKeyError := jwt.ParseRSAPublicKeyFromPEM(publicKeyContents)
	if decodedPublicKeyError != nil {
		log.Fatal("Error parsing public key: ", decodedPublicKeyError)
	}

	_privateKey = decodedPrivateKey
	_publicKey = decodedPublicKey
}

// Sign will generate new jwt token.
func Sign(userId string, username string) (tokenString string, err error) {
	// create claims
	var jwtPayload = make(jwt.MapClaims)
	jwtPayload["aud"] = "DEVORIA"
	jwtPayload["exp"] = time.Now().Add(time.Hour * 24).Unix()
	jwtPayload["iat"] = time.Now().Unix()
	jwtPayload["iss"] = "DEVORIA"
	jwtPayload["jti"] = userId
	jwtPayload["username"] = username

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtPayload)

	// return signed token
	return token.SignedString(_privateKey)
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := apm.StartSpan(c.Request.Context(), "JwtAuthMiddleware", "http")
		defer span.End()

		// get the authorization header from the request
		authorizationHeader := c.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, "Missing authorization header", nil))
			return
		}

		// check if the authentication method is Bearer
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, "Authorization header is not Bearer", nil))
			return
		}

		// remove Bearer
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		// parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate whether the token is signed using RSA
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("invalid token")
			}

			// Token validation success, return public key to validate with
			return _publicKey, nil
		})

		// check if the decoding process has any errors
		if err != nil {
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, err.Error(), nil))
			return
		}

		// check if the token has validation errors
		jwtErrors, _ := err.(*jwt.ValidationError)
		if jwtErrors != nil && jwtErrors.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, "Token is expired or not valid yet", nil))
		}

		// decode token
		mapClaims := token.Claims.(jwt.MapClaims)

		var jwtPayload JwtPayload
		jwtPayload.Audience = mapClaims["aud"].(string)
		jwtPayload.ExpiresAt = int64(mapClaims["exp"].(float64))
		jwtPayload.IssuedAt = int64(mapClaims["iat"].(float64))
		jwtPayload.Issuer = mapClaims["iss"].(string)
		jwtPayload.Id = mapClaims["jti"].(string)
		jwtPayload.Username = mapClaims["username"].(string)

		// set token to the context
		c.Set("JWT_AUTHENTICATED", true)
		c.Set("JWT_PAYLOAD", &jwtPayload)

		// resume chain
		c.Next()
	}
}

func GetJwtUser(c *gin.Context) (user *JwtPayload, err error) {
	// check if the token is set
	if _, exists := c.Get("JWT_AUTHENTICATED"); !exists {
		return nil, fmt.Errorf("JWT_AUTHENTICATED not found in context")
	}

	// get the token from the context
	return c.MustGet("JWT_PAYLOAD").(*JwtPayload), nil
}
