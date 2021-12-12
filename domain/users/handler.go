package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

type UserHandler struct {
	Usecase UserUseCase
}

func ConstructUserHandler(router *gin.Engine, usecase UserUseCase, jwtAuth authentication.IJwtAuth) {
	handler := &UserHandler{
		Usecase: usecase,
	}

	v1 := router.Group("/api/v1/user")
	v1.POST("/login", authentication.BasicAuthMiddleware(), handler.Login)
	v1.POST("/register", authentication.BasicAuthMiddleware(), handler.Register)
	v1.GET("/me", jwtAuth.JwtAuthMiddleware(), handler.Profile)
}

// login godoc
// @Summary      Login
// @Description  Authenticate a user to get access token
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body   string true "Login body"
// @Router       /api/v1/user/login [post]
// @Security     BasicAuth
func (u *UserHandler) Login(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Login", "http")
	defer span.End()

	var model LoginModel
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	result := u.Usecase.Login(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Register new account
// @Description  Register new account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body   string true "Register body"
// @Router       /api/v1/user/register [post]
// @Security     BasicAuth
func (u *UserHandler) Register(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Register", "http")
	defer span.End()

	var model RegisterModel
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	result := u.Usecase.Register(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Get profile
// @Description  Get detailed user profile
// @Tags         accounts
// @Produce      json
// @Router       /api/v1/user/me [get]
// @Security     JwtAuth
func (u *UserHandler) Profile(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Profile", "http")
	defer span.End()

	user, err := authentication.GetJwtUser(c)
	if err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusInternalServerError, "Can't get user from token", nil))
		return
	}

	userId, _ := strconv.Atoi(user.Id)
	model := GetProfileModel{
		ID: userId,
	}

	result := u.Usecase.GetProfile(c, model)
	utils.WriteResponse(c, result)
}
