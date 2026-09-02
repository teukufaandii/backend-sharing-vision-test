package routes

import (
	"strings"

	handler "backend-sharing-vision-test/internal/handler/http"
	"backend-sharing-vision-test/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter initializes and configures the Gin router, middlewares, and routes
func SetupRouter(allowedOrigins []string, articleHandler *handler.ArticleHandler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORSMiddleware(allowedOrigins))
	r.Use(middleware.LoggerMiddleware(logrus.New()))
	r.Use(middleware.RecoveryMiddleware(logrus.New()))

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	if articleHandler != nil {
		// 1. POST /article/ and POST /article
		r.POST("/article", articleHandler.CreateArticle)
		r.POST("/article/", articleHandler.CreateArticle)

		// 2. GET /article/<limit>/<offset>
		r.GET("/article/:first/:second", articleHandler.GetArticles)

		// 3. GET /article/<id>
		r.GET("/article/:first", articleHandler.GetArticleByID)

		// 4. PUT / PATCH /article/<id> (Update)
		r.PUT("/article/:first", articleHandler.UpdateArticle)
		r.PATCH("/article/:first", articleHandler.UpdateArticle)

		// 5. DELETE /article/<id> (Delete)
		r.DELETE("/article/:first", articleHandler.DeleteArticle)

		// POST /article/<id>
		r.POST("/article/:first", func(c *gin.Context) {
			action := c.Query("action")
			methodOverride := c.GetHeader("X-HTTP-Method-Override")

			if strings.EqualFold(action, "delete") || strings.EqualFold(methodOverride, "DELETE") {
				articleHandler.DeleteArticle(c)
				return
			}
			articleHandler.UpdateArticle(c)
		})
	}

	return r
}
