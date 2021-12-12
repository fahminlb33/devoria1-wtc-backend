package articles

import "time"

type ArticleDto struct {
	ID        uint                 `json:"id"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	Slug      string               `json:"slug"`
	Status    ArticlePublishStatus `json:"status"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

type ArticleItemDto struct {
	ID        uint                 `json:"id"`
	Title     string               `json:"title"`
	Slug      string               `json:"slug"`
	Status    ArticlePublishStatus `json:"status"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}
