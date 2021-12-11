package articles

type FindAllModel struct {
	UserId  int
	Keyword string `form:"email"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

type GetModel struct {
	UserId    int
	ArticleId int `uri:"id" binding:"required"`
}

type CreateModel struct {
	UserId  int
	Title   string               `json:"title" binding:"required"`
	Content string               `json:"content" binding:"required"`
	Slug    string               `json:"slug" binding:"required"`
	Status  ArticlePublishStatus `json:"status" binding:"required"`
	IsDraft bool                 `json:"isDraft" binding:"required"`
}

type SaveModel struct {
	UserId    int
	ArticleId int                  `json:"id" binding:"required"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	Slug      string               `json:"slug"`
	Status    ArticlePublishStatus `json:"status"`
	IsDraft   bool                 `json:"isDraft"`
}

type DeleteModel struct {
	UserId    int
	ArticleId int `uri:"id" binding:"required"`
}
