package main

import (
	"log"
	"os"

	"backend-sharing-vision-test/config"
	"backend-sharing-vision-test/internal/routes"
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
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		// Auto-migrate schema on start
		if err := dbMysql.AutoMigrate(db); err != nil {
			log.Printf("Warning: Auto-migration failed: %v", err)
		}
	}

	// Setup HTTP router
	router := routes.SetupRouter(cfg.AllowedOrigins)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
