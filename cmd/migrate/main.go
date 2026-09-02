package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"backend-sharing-vision-test/config"
	dbMysql "backend-sharing-vision-test/pkg/database/mysql"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	cfg := config.Load()

	db, err := dbMysql.NewMySQLConn(cfg.DatabaseURL, cfg.DatabaseMaxConn, cfg.DatabaseMaxIdle)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	migrationsPath := "database/migrations"
	if customPath := os.Getenv("MIGRATIONS_PATH"); customPath != "" {
		migrationsPath = customPath
	}

	var action string
	if len(os.Args) > 1 {
		action = os.Args[1]
	} else {
		flag.StringVar(&action, "action", "up", "Migration action: up, down, or auto")
		flag.Parse()
	}

	switch action {
	case "up":
		log.Println("Executing SQL migrations (UP)...")
		if err := dbMysql.RunSQLMigrations(sqlDB, migrationsPath); err != nil {
			log.Fatalf("Migration UP failed: %v", err)
		}
		log.Println("Migration UP completed successfully.")

	case "down":
		log.Println("Executing SQL migrations (DOWN)...")
		if err := dbMysql.RollbackSQLMigrations(sqlDB, migrationsPath, 1); err != nil {
			log.Fatalf("Migration DOWN failed: %v", err)
		}
		log.Println("Migration DOWN completed successfully.")

	case "down-all":
		log.Println("Executing all SQL migrations rollback (DOWN ALL)...")
		if err := dbMysql.RollbackSQLMigrations(sqlDB, migrationsPath, 0); err != nil {
			log.Fatalf("Migration DOWN ALL failed: %v", err)
		}
		log.Println("Migration DOWN ALL completed successfully.")

	case "auto":
		log.Println("Executing GORM AutoMigrate...")
		if err := dbMysql.AutoMigrate(db); err != nil {
			log.Fatalf("GORM AutoMigrate failed: %v", err)
		}
		log.Println("GORM AutoMigrate completed successfully.")

	default:
		fmt.Printf("Unknown migration action: %s\nUsage: go run cmd/migrate/main.go [up|down|down-all|auto]\n", action)
		os.Exit(1)
	}
}
