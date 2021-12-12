package articles

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.elastic.co/apm"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

type ArticleHandler struct {
	Validator *validator.Validate
	Usecase   ArticleUseCase
}

func ConstructArticlesHandler(router *gin.Engine, usecase ArticleUseCase) {
	handler := &ArticleHandler{
		Usecase: usecase,
	}

	v1 := router.Group("/api/v1/articles")
	v1.GET("", authentication.JwtAuthMiddleware(), handler.FindAll)
	v1.POST("", authentication.JwtAuthMiddleware(), handler.Create)
	v1.PUT("", authentication.JwtAuthMiddleware(), handler.Save)
	v1.GET("/:id", authentication.JwtAuthMiddleware(), handler.Get)
	v1.DELETE("/:id", authentication.JwtAuthMiddleware(), handler.Delete)
}

// @Summary      Create article
// @Description  Create single article
// @Tags         articles
// @Produce      json
// @Param        keyword   path      string  false  "Keyword in title or content"
// @Param 	     page      path      int     false  "Pagination number"
// @Param 	     limit     path      int     false  "Number of rows to be retrieved"
// @Router       /api/v1/articles/ [get]
func (u *ArticleHandler) FindAll(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "FindAll", "http")
	defer span.End()

	var model FindAllModel
	if err := c.ShouldBindQuery(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	injectUserId(c, &model.UserId)

	result := u.Usecase.FindAll(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Create article
// @Description  Create single article
// @Tags         articles
// @Accepts      json
// @Produce      json
// @Router       /api/v1/articles [post]
func (u *ArticleHandler) Create(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Create", "http")
	defer span.End()

	var model CreateModel
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	injectUserId(c, &model.UserId)

	result := u.Usecase.Create(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Save article
// @Description  Update single article
// @Tags         articles
// @Accepts      json
// @Produce      json
// @Router       /api/v1/articles [put]
func (u *ArticleHandler) Save(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Save", "http")
	defer span.End()

	var model SaveModel
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	injectUserId(c, &model.UserId)

	result := u.Usecase.Save(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Show article
// @Description  Get single article
// @Tags         articles
// @Produce      json
// @Param        id   path      int  true  "Article ID"
// @Router       /api/v1/articles/:id [get]
func (u *ArticleHandler) Get(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Get", "http")
	defer span.End()

	var model GetModel
	if err := c.ShouldBindUri(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	injectUserId(c, &model.UserId)

	result := u.Usecase.Get(c, model)
	utils.WriteResponse(c, result)
}

// @Summary      Delete an article
// @Description  Delete single article
// @Tags         articles
// @Produce      json
// @Param        id   path      int  true  "Article ID"
// @Router       /api/v1/articles/:id [delete]
func (u *ArticleHandler) Delete(c *gin.Context) {
	span, _ := apm.StartSpan(c.Request.Context(), "Delete", "http")
	defer span.End()

	var model DeleteModel
	if err := c.ShouldBindUri(&model); err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusBadRequest, "Validation failed", err.Error()))
		return
	}

	injectUserId(c, &model.UserId)

	result := u.Usecase.Delete(c, model)
	utils.WriteResponse(c, result)
}

func injectUserId(c *gin.Context, userId *int) {
	user, err := authentication.GetJwtUser(c)
	if err != nil {
		utils.WriteResponse(c, utils.WrapResponse(http.StatusInternalServerError, "Can't get user from token", nil))
		return
	}

	id, _ := strconv.Atoi(user.Id)
	*userId = id
}
