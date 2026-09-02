package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handler "backend-sharing-vision-test/internal/handler/http"
	"backend-sharing-vision-test/internal/models"
	"backend-sharing-vision-test/internal/repository"
	"backend-sharing-vision-test/internal/routes"
	"backend-sharing-vision-test/internal/service"

	"github.com/gin-gonic/gin"
)

type mockRepo struct {
	posts  map[int]*models.Post
	lastID int
}

func newMockRepo() *mockRepo {
	return &mockRepo{posts: make(map[int]*models.Post)}
}

func (m *mockRepo) Create(ctx context.Context, post *models.Post) error {
	m.lastID++
	post.ID = m.lastID
	post.CreatedDate = time.Now()
	post.UpdatedDate = time.Now()
	stored := *post
	m.posts[post.ID] = &stored
	return nil
}

func (m *mockRepo) FindAll(ctx context.Context, limit, offset int) ([]models.Post, error) {
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

func (m *mockRepo) FindByID(ctx context.Context, id int) (*models.Post, error) {
	if p, ok := m.posts[id]; ok {
		copyP := *p
		return &copyP, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockRepo) Update(ctx context.Context, post *models.Post) error {
	if _, ok := m.posts[post.ID]; !ok {
		return repository.ErrNotFound
	}
	post.UpdatedDate = time.Now()
	stored := *post
	m.posts[post.ID] = &stored
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	if _, ok := m.posts[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.posts, id)
	return nil
}

func (m *mockRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(m.posts)), nil
}

func setupTestRouter() (*gin.Engine, *mockRepo) {
	gin.SetMode(gin.TestMode)
	repo := newMockRepo()
	svc := service.NewArticleService(repo)
	h := handler.NewArticleHandler(svc)
	r := routes.SetupRouter([]string{"*"}, h)
	return r, repo
}

func TestArticleEndpoints(t *testing.T) {
	validContent := strings.Repeat("This is a valid long content body for our testing purpose. ", 8) // > 200 chars

	t.Run("POST /article/ - Success (201)", func(t *testing.T) {
		router, _ := setupTestRouter()

		payload := map[string]string{
			"title":    "Valid Article Title That Has More Than 20 Characters",
			"content":  validContent,
			"category": "Technology",
			"status":   "Publish",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/article/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST /article/ - Validation Error (400)", func(t *testing.T) {
		router, _ := setupTestRouter()

		payload := map[string]string{
			"title":    "Short Title",
			"content":  "Short",
			"category": "IT",
			"status":   "invalid-status",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/article/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET /article/<limit>/<offset> - Success (200)", func(t *testing.T) {
		router, repo := setupTestRouter()

		for i := 1; i <= 5; i++ {
			_ = repo.Create(context.Background(), &models.Post{
				Title:    "Article Title For Pagination Item Number " + string(rune('0'+i)),
				Content:  validContent,
				Category: "Tech",
				Status:   "Draft",
			})
		}

		req, _ := http.NewRequest(http.MethodGet, "/article/2/0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var items []map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &items)
		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("GET /article/<id> - Found (200) and Not Found (404)", func(t *testing.T) {
		router, repo := setupTestRouter()

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Specific Article Title For ID Query Test",
			Content:  validContent,
			Category: "Tech",
			Status:   "Draft",
		})

		// 200 OK
		req, _ := http.NewRequest(http.MethodGet, "/article/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		// 404 Not Found
		req404, _ := http.NewRequest(http.MethodGet, "/article/999", nil)
		w404 := httptest.NewRecorder()
		router.ServeHTTP(w404, req404)

		if w404.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w404.Code)
		}
	})

	t.Run("PUT /article/<id> - Success (200)", func(t *testing.T) {
		router, repo := setupTestRouter()

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Original Title Long Enough For Validation",
			Content:  validContent,
			Category: "Tech",
			Status:   "Draft",
		})

		payload := map[string]string{
			"title":    "Updated Title Long Enough For Validation Test",
			"content":  validContent + " with extra modifications.",
			"category": "Updated Category",
			"status":   "Publish",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/article/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DELETE /article/<id> - Success (200) and Not Found (404)", func(t *testing.T) {
		router, repo := setupTestRouter()

		_ = repo.Create(context.Background(), &models.Post{
			Title:    "Article To Delete Long Enough For Testing",
			Content:  validContent,
			Category: "Tech",
			Status:   "Thrash",
		})

		req, _ := http.NewRequest(http.MethodDelete, "/article/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		req2, _ := http.NewRequest(http.MethodDelete, "/article/1", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w2.Code)
		}
	})
}
