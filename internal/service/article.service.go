package service

import (
	"context"
	"fmt"

	"backend-sharing-vision-test/internal/dto/request"
	"backend-sharing-vision-test/internal/models"
	"backend-sharing-vision-test/internal/repository"
	"backend-sharing-vision-test/pkg/utils"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, req *request.CreateArticleRequest) (*models.Post, error)
	GetArticles(ctx context.Context, limit, offset int) ([]models.Post, error)
	GetArticleByID(ctx context.Context, id int) (*models.Post, error)
	UpdateArticle(ctx context.Context, id int, req *request.UpdateArticleRequest) (*models.Post, error)
	DeleteArticle(ctx context.Context, id int) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) CreateArticle(ctx context.Context, req *request.CreateArticleRequest) (*models.Post, error) {
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	post := &models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *articleService) GetArticles(ctx context.Context, limit, offset int) ([]models.Post, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *articleService) GetArticleByID(ctx context.Context, id int) (*models.Post, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *articleService) UpdateArticle(ctx context.Context, id int, req *request.UpdateArticleRequest) (*models.Post, error) {
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Title = req.Title
	existing.Content = req.Content
	existing.Category = req.Category
	existing.Status = req.Status

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
