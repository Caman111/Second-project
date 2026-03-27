package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `gorm:"not null" validate:"required,min=3"`
	Description string         `validate:"required"`
	Images      pq.StringArray `gorm:"type:text[]"`
}
