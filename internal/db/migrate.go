package db

import (
	"gorm.io/gorm"
	"zl0y.team/billing/internal/models"
)

func MigratePostgres(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
