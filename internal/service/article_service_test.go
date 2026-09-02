package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"backend-sharing-vision-test/internal/dto/request"
	"backend-sharing-vision-test/internal/models"
	"backend-sharing-vision-test/internal/repository"
	"backend-sharing-vision-test/internal/service"
)

type MockArticleRepository struct {
	posts  map[int]*models.Post
	lastID int
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		posts: make(map[int]*models.Post),
	}
}

func (m *MockArticleRepository) Create(ctx context.Context, post *models.Post) error {
	m.lastID++
	post.ID = m.lastID
	post.CreatedDate = time.Now()
	post.UpdatedDate = time.Now()
	stored := *post
	m.posts[post.ID] = &stored
	return nil
}

func (m *MockArticleRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Post, error) {
	var list []models.Post
	for i := 1; i <= m.lastID; i++ {
		if p, ok := m.posts[i]; ok {
			list = append(list, *p)
		}
	}

	if offset >= len(list) {
		return []models.Post{}, nil
	}

	end := len(list)
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	return list[offset:end], nil
}

func (m *MockArticleRepository) FindByID(ctx context.Context, id int) (*models.Post, error) {
	if p, ok := m.posts[id]; ok {
		copyPost := *p
		return &copyPost, nil
	}
	return nil, repository.ErrNotFound
}

func (m *MockArticleRepository) Update(ctx context.Context, post *models.Post) error {
	if _, ok := m.posts[post.ID]; !ok {
		return repository.ErrNotFound
	}
	post.UpdatedDate = time.Now()
	stored := *post
	m.posts[post.ID] = &stored
	return nil
}

func (m *MockArticleRepository) Delete(ctx context.Context, id int) error {
	if _, ok := m.posts[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.posts, id)
	return nil
}

func (m *MockArticleRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.posts)), nil
}

func TestArticleService(t *testing.T) {
	validContent := strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", 8) // > 200 chars

	t.Run("CreateArticle - Success", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		req := &request.CreateArticleRequest{
			Title:    "Valid Article Title That Has More Than 20 Characters",
			Content:  validContent,
			Category: "Technology",
			Status:   "Publish",
		}

		post, err := svc.CreateArticle(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if post.ID == 0 {
			t.Errorf("expected generated ID > 0, got 0")
		}
		if post.Title != req.Title {
			t.Errorf("expected title '%s', got '%s'", req.Title, post.Title)
		}
	})

	t.Run("CreateArticle - Validation Error", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		req := &request.CreateArticleRequest{
			Title:    "Too Short",
			Content:  validContent,
			Category: "Technology",
			Status:   "Publish",
		}

		_, err := svc.CreateArticle(context.Background(), req)
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("GetArticles - Paging", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		for i := 1; i <= 5; i++ {
			_ = repo.Create(context.Background(), &models.Post{
				Title:    "Article Title For Item Number " + string(rune('0'+i)),
				Content:  validContent,
				Category: "News",
				Status:   "Publish",
			})
		}

		posts, err := svc.GetArticles(context.Background(), 2, 0)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("expected 2 items, got %d", len(posts))
		}

		postsOffset, err := svc.GetArticles(context.Background(), 2, 2)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(postsOffset) != 2 {
			t.Errorf("expected 2 items, got %d", len(postsOffset))
		}
		if postsOffset[0].ID != 3 {
			t.Errorf("expected first item ID 3, got %d", postsOffset[0].ID)
		}
	})

	t.Run("GetArticleByID - Found and Not Found", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Single Article Title Long Enough Here",
			Content:  validContent,
			Category: "News",
			Status:   "Draft",
		})

		post, err := svc.GetArticleByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if post.ID != 1 {
			t.Errorf("expected ID 1, got %d", post.ID)
		}

		_, err = svc.GetArticleByID(context.Background(), 999)
		if err == nil {
			t.Fatalf("expected ErrNotFound, got nil")
		}
	})

	t.Run("UpdateArticle - Success and Not Found", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Initial Article Title With Enough Length",
			Content:  validContent,
			Category: "News",
			Status:   "Draft",
		})

		updateReq := &request.UpdateArticleRequest{
			Title:    "Updated Article Title With Enough Length",
			Content:  validContent + " updated details",
			Category: "Updated Category",
			Status:   "Publish",
		}

		updated, err := svc.UpdateArticle(context.Background(), 1, updateReq)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if updated.Title != updateReq.Title {
			t.Errorf("expected updated title '%s', got '%s'", updateReq.Title, updated.Title)
		}

		_, err = svc.UpdateArticle(context.Background(), 999, updateReq)
		if err == nil {
			t.Fatalf("expected error for non-existent ID, got nil")
		}
	})

	t.Run("DeleteArticle - Success and Not Found", func(t *testing.T) {
		repo := NewMockArticleRepository()
		svc := service.NewArticleService(repo)

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Article To Be Deleted Long Title",
			Content:  validContent,
			Category: "News",
			Status:   "Thrash",
		})

		err := svc.DeleteArticle(context.Background(), 1)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		err = svc.DeleteArticle(context.Background(), 1)
		if err == nil {
			t.Fatalf("expected error deleting non-existent article, got nil")
		}
	})
}
