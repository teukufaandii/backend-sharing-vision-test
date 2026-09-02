package routes

import (
	handler "backend-sharing-vision-test/internal/handler/http"
	"backend-sharing-vision-test/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter initializes and configures the Gin router and middlewares
func SetupRouter(allowedOrigins []string) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORSMiddleware(allowedOrigins))
	r.Use(middleware.LoggerMiddleware(logrus.New()))
	r.Use(middleware.RecoveryMiddleware(logrus.New()))

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	return r
}
