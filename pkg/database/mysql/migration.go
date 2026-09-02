package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"backend-sharing-vision-test/internal/models"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// AutoMigrate performs schema migration using GORM AutoMigrate
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting GORM AutoMigrate...")

	if err := db.AutoMigrate(&models.Post{}); err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	log.Println("GORM AutoMigrate completed successfully")
	return nil
}

// RunSQLMigrations runs up migrations from the specified migrations path using golang-migrate
func RunSQLMigrations(sqlDB *sql.DB, migrationsPath string) error {
	log.Printf("Running SQL migrations from '%s'...", migrationsPath)

	driver, err := migratemysql.WithInstance(sqlDB, &migratemysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create MySQL migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("SQL migrations applied successfully")
	return nil
}

// RollbackSQLMigrations runs down migrations (1 step or all) using golang-migrate
func RollbackSQLMigrations(sqlDB *sql.DB, migrationsPath string, steps int) error {
	log.Printf("Rolling back SQL migrations from '%s' (steps: %d)...", migrationsPath, steps)

	driver, err := migratemysql.WithInstance(sqlDB, &migratemysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create MySQL migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration instance: %w", err)
	}

	var errRollback error
	if steps > 0 {
		errRollback = m.Steps(-steps)
	} else {
		errRollback = m.Down()
	}

	if errRollback != nil && !errors.Is(errRollback, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migrations: %w", errRollback)
	}

	log.Println("SQL rollback completed successfully")
	return nil
}
