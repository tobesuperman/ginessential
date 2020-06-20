package vo

type CreatePostRequest struct {
	CategoryId uint   `json:"categoryId" binding:"required"`
	Title      string `json:"title" binding:"required,max=10"`
	HeadImg    string `json:"headImg"`
	Content    string `json:"content" binding:"required"`
}
