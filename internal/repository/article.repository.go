package repository

import (
	"context"
	"errors"

	"backend-sharing-vision-test/internal/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("article not found")
)

type ArticleRepository interface {
	Create(ctx context.Context, post *models.Post) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Post, error)
	FindByID(ctx context.Context, id int) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int64, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *articleRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	query := r.db.WithContext(ctx).Order("id ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&posts).Error
	return posts, err
}

func (r *articleRepository) FindByID(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *articleRepository) Update(ctx context.Context, post *models.Post) error {
	result := r.db.WithContext(ctx).Model(&models.Post{}).
		Where("id = ?", post.ID).
		Updates(map[string]interface{}{
			"title":    post.Title,
			"content":  post.Content,
			"category": post.Category,
			"status":   post.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Post{}).Count(&count).Error
	return count, err
}
