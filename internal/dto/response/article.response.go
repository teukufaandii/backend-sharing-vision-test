package response

import (
	"time"

	"backend-sharing-vision-test/internal/models"
)

type ArticleResponse struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
	Status      string    `json:"status"`
}

func FromPostModel(post *models.Post) ArticleResponse {
	if post == nil {
		return ArticleResponse{}
	}
	return ArticleResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Category:    post.Category,
		CreatedDate: post.CreatedDate,
		UpdatedDate: post.UpdatedDate,
		Status:      post.Status,
	}
}

func FromPostModelList(posts []models.Post) []ArticleResponse {
	responses := make([]ArticleResponse, 0, len(posts))
	for i := range posts {
		responses = append(responses, FromPostModel(&posts[i]))
	}
	return responses
}
