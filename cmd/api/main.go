package main

import (
	"log"
	"os"

	"backend-sharing-vision-test/config"
	handler "backend-sharing-vision-test/internal/handler/http"
	"backend-sharing-vision-test/internal/repository"
	"backend-sharing-vision-test/internal/routes"
	"backend-sharing-vision-test/internal/service"
	dbMysql "backend-sharing-vision-test/pkg/database/mysql"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

	// Initialize MySQL database connection
	db, err := dbMysql.NewMySQLConn(cfg.DatabaseURL, cfg.DatabaseMaxConn, cfg.DatabaseMaxIdle)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate schema on start
	if err := dbMysql.AutoMigrate(db); err != nil {
		log.Printf("Warning: Auto-migration failed: %v", err)
	}

	// Initialize repository, service, and handler layers
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	// Setup HTTP router
	router := routes.SetupRouter(cfg.AllowedOrigins, articleHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
