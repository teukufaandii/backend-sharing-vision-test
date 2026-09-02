package request

type CreateArticleRequest struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash Publish Draft Thrash"`
}

type UpdateArticleRequest struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash Publish Draft Thrash"`
}
