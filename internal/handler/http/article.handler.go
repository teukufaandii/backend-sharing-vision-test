package handler

import (
	"errors"
	"net/http"
	"strconv"

	"backend-sharing-vision-test/internal/dto/request"
	"backend-sharing-vision-test/internal/dto/response"
	"backend-sharing-vision-test/internal/repository"
	"backend-sharing-vision-test/internal/service"
	"backend-sharing-vision-test/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ArticleHandler handles HTTP requests for article resource
type ArticleHandler struct {
	service service.ArticleService
}

// NewArticleHandler creates a new instance of ArticleHandler
func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// CreateArticle handles POST /article/
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := utils.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := h.service.CreateArticle(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetArticles handles GET /article/:limit/:offset (or query params limit & offset)
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	limitStr := c.Param("limit")
	if limitStr == "" {
		limitStr = c.Param("first")
	}
	offsetStr := c.Param("offset")
	if offsetStr == "" {
		offsetStr = c.Param("second")
	}

	if limitStr == "" {
		limitStr = c.DefaultQuery("limit", "10")
	}
	if offsetStr == "" {
		offsetStr = c.DefaultQuery("offset", "0")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be a non-negative integer",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "offset must be a non-negative integer",
		})
		return
	}

	posts, err := h.service.GetArticles(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := response.FromPostModelList(posts)
	c.JSON(http.StatusOK, res)
}

// GetArticleByID handles GET /article/:id
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Param("first")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must be a valid positive integer",
		})
		return
	}

	post, err := h.service.GetArticleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := response.FromPostModel(post)
	c.JSON(http.StatusOK, res)
}

// UpdateArticle handles PUT/PATCH/POST /article/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Param("first")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must be a valid positive integer",
		})
		return
	}

	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := utils.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = h.service.UpdateArticle(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteArticle handles DELETE /article/:id (and POST /article/:id/delete or DELETE action)
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Param("first")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must be a valid positive integer",
		})
		return
	}

	err = h.service.DeleteArticle(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
