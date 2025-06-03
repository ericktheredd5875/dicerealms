package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

var DB *gorm.DB

func InitDB() error {

	dsn := buildDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&PlayerModel{},
		&RoomModel{},
		&ItemModel{},
		&SceneModel{},
		&SceneLogModel{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	log.Println("Database schema migrated")
	return nil
}

func buildDSN() string {

	dbHost := utils.RequiredEnv("DB_HOST")
	dbPort := utils.RequiredEnv("DB_PORT")
	dbUser := utils.RequiredEnv("DB_USER")
	dbPass := utils.RequiredEnv("DB_PASS")
	dbName := utils.RequiredEnv("DB_NAME")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
}
