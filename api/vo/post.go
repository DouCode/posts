package vo

type CreatePostRequest struct {
	CategoryId string `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=10"`
	HeadImg    string `json:"head_img"`
	Content    string `json:"content" binding:"required"`
}

type CreateBlogRequest struct {
	Title   string `json:"blogTitle" binding:"required,max=10"`
	Content string `json:"blogBody" binding:"required"`
	TagStr  string `json:"tagId" binding:"required"`
}
