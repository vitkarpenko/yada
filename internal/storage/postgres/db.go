package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"yada/internal/config"
)

func NewDB(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Reminder{}); err != nil {
		return nil, err
	}

	return db, nil
}
