package authentication

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := apm.StartSpan(c.Request.Context(), "BasicAuthMiddleware", "http")
		defer span.End()

		// get the authorization header from the request
		authorizationHeader := c.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			c.Header("WWW-Authenticate", "Basic realm=DEVORIA")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}

		// check if the authentication method is Basic
		if !strings.HasPrefix(authorizationHeader, "Basic ") {
			c.Header("WWW-Authenticate", "Basic realm=DEVORIA")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is not Bearer"})
			return
		}

		// remove Basic
		tokenString := strings.Replace(authorizationHeader, "Basic ", "", -1)

		// parse token
		tokenBody, err := base64.StdEncoding.DecodeString(tokenString)

		// check if the decoding process has any errors
		if err != nil {
			c.Header("WWW-Authenticate", "Basic realm=DEVORIA")
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, err.Error(), nil))
			return
		}

		// check if the token has validation errors
		parts := strings.Split(string(tokenBody), ":")

		// verify credentials
		if !SafeCompareString(parts[0], config.GlobalConfig.Authentication.BasicUsername) || !SafeCompareString(parts[1], config.GlobalConfig.Authentication.BasicPassword) {
			c.Header("WWW-Authenticate", "Basic realm=DEVORIA")
			utils.WriteAbortResponse(c, utils.WrapResponse(http.StatusUnauthorized, "Invalid credentials", nil))
			return
		}

		// set token to the context
		c.Set("BASIC_AUTHENTICATED", true)

		// resume chain
		c.Next()
	}
}
