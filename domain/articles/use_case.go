package articles

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

type ArticleUseCase interface {
	FindAll(c context.Context, model FindAllModel) (response utils.Response)
	Get(c context.Context, model GetModel) (response utils.Response)
	Create(c context.Context, model CreateModel) (response utils.Response)
	Save(c context.Context, model SaveModel) (response utils.Response)
	Delete(c context.Context, model DeleteModel) (response utils.Response)
}

type ArticleUseCaseImpl struct {
	Database *gorm.DB
}

func ConstructArticlesUseCase(db *gorm.DB) ArticleUseCase {
	return &ArticleUseCaseImpl{
		Database: db,
	}
}

func (u *ArticleUseCaseImpl) FindAll(c context.Context, model FindAllModel) (response utils.Response) {
	db := u.Database.WithContext(c)

	// get the author
	var author users.User
	db.First(&author, model.UserId)

	// get the article
	var articles []Article
	db.Preload(clause.Associations).Where("author_id = ?", model.UserId).Find(&articles)

	// create response
	finalResponse := []ArticleItemDto{}
	for _, article := range articles {
		finalResponse = append(finalResponse, ArticleItemDto{
			Title:  article.Title,
			Slug:   article.Slug,
			Status: article.Status,
		})
	}

	return utils.WrapResponse(http.StatusOK, "OK", finalResponse)
}

func (u *ArticleUseCaseImpl) Get(c context.Context, model GetModel) (response utils.Response) {
	db := u.Database.WithContext(c)

	// get the author
	var author users.User
	db.First(&author, model.UserId)

	// get the article
	var article Article
	db.Preload(clause.Associations).First(&article, model.ArticleId)

	// check whether the article is owned by the user
	if author.ID != article.Author.ID {
		return utils.WrapResponse(http.StatusOK, "Article not found", nil)
	}

	// create response
	finalResponse := ArticleDto{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Slug:      article.Slug,
		Status:    article.Status,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	return utils.WrapResponse(http.StatusOK, "OK", finalResponse)
}

func (u *ArticleUseCaseImpl) Create(c context.Context, model CreateModel) (response utils.Response) {
	db := u.Database.WithContext(c)

	// get the author
	var author users.User
	db.First(&author, model.UserId)

	// create entity
	article := Article{
		Title:   model.Title,
		Content: model.Content,
		Slug:    model.Slug,
		Status:  model.Status,
		Author:  author,
	}

	db.Create(&article)

	// create response
	finalResponse := ArticleDto{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Slug:      article.Slug,
		Status:    article.Status,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	return utils.WrapResponse(http.StatusCreated, "Article created", finalResponse)
}

func (u *ArticleUseCaseImpl) Save(c context.Context, model SaveModel) (response utils.Response) {
	return
}

func (u *ArticleUseCaseImpl) Delete(c context.Context, model DeleteModel) (response utils.Response) {
	db := u.Database.WithContext(c)

	// get the author
	var author users.User
	db.First(&author, model.UserId)

	// get the article
	var article Article
	db.Preload(clause.Associations).First(&article, model.ArticleId)

	// check whether the article is owned by the user
	if author.ID != article.Author.ID {
		return utils.WrapResponse(http.StatusOK, "Article not found", nil)
	}

	// delete article
	db.Delete(&article)

	return utils.WrapResponse(http.StatusOK, "OK", nil)
}