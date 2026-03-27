package db

import (
	models "3-validation-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Product{}, &models.BiznessCreate{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
