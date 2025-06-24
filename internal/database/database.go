package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	"WeatherAPITask/internal/config"
	"WeatherAPITask/internal/models"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.Weather{}); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
