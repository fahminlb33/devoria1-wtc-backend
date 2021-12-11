package articles

import (
	"gorm.io/gorm"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
)

type ArticlePublishStatus string

const (
	DRAFT     ArticlePublishStatus = "DRAFT"
	PUBLISHED ArticlePublishStatus = "PUBLISHED"
	ARCHIVED  ArticlePublishStatus = "ARCHIVED"
)

type Article struct {
	gorm.Model
	Title   string               `gorm:"size:255;not null;index"`
	Content string               `gorm:"type:text;not null"`
	Slug    string               `gorm:"size:255;not null,uniqueIndex"`
	Status  ArticlePublishStatus `gorm:"size:255;not null"`

	AuthorId int
	Author   users.User `gorm:"foreignKey:AuthorId"`
}
