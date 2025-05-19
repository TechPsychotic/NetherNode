package storage

import (
	"fmt"
	"your-project/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(config *utils.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Server{})

	return db
}