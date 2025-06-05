package db

import (
	"database/sql"

	"gorm.io/driver/postgres" // or postgres, sqlite, etc., based on your actual DB
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormFromSQLDB(sqlDB *sql.DB) (*gorm.DB, error) {
	// This assumes you're using a compatible driver with sqlmock (most often `mysql` or `postgres`)
	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	// Turn off logging during tests for cleaner output
	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
