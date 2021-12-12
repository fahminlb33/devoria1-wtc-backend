package articles

type FindAllModel struct {
	UserId  int
	Keyword string `form:"keyword"`
	Page    int    `form:"page" binding:"omitempty,numeric"`
	Limit   int    `form:"limit" binding:"omitempty,numeric"`
}

type CreateModel struct {
	UserId  int
	Title   string               `json:"title" binding:"required,min=5,max=255"`
	Content string               `json:"content" binding:"required,min=10"`
	Slug    string               `json:"slug" binding:"required,min=5,max=255"`
	Status  ArticlePublishStatus `json:"status" binding:"required,isarticlepublishstatus"`
}

type SaveModel struct {
	UserId    int
	ArticleId int                  `json:"id" binding:"required,numeric"`
	Title     string               `json:"title" binding:"omitempty,min=5,max=255"`
	Content   string               `json:"content" binding:"omitempty,min=10"`
	Slug      string               `json:"slug" binding:"omitempty,min=5,max=255"`
	Status    ArticlePublishStatus `json:"status" binding:"omitempty,isarticlepublishstatus"`
}

type GetModel struct {
	UserId    int
	ArticleId int `uri:"id" binding:"required,numeric"`
}

type DeleteModel struct {
	UserId    int
	ArticleId int `uri:"id" binding:"required,numeric"`
}
