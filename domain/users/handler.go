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

func ConstructUserHandler(router *gin.Engine, usecase UserUseCase) {
	handler := &UserHandler{
		Usecase: usecase,
	}

	v1 := router.Group("/api/v1/user")
	v1.POST("/login", authentication.BasicAuthMiddleware(), handler.Login)
	v1.POST("/register", authentication.BasicAuthMiddleware(), handler.Register)
	v1.GET("/me", authentication.JwtAuthMiddleware(), handler.Profile)
}

// login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /api/v1/user/login [post]
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

// login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /api/v1/user/register [post]
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

// login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /api/v1/user/me [get]
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
