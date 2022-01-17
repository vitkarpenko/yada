package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"yada/internal/config"
)

type DB struct {
	*gorm.DB
}

func NewDB(config config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Reminder{}, &LastQuote{}); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (d *DB) GetLastQuoteHash() string {
	lastQuote := LastQuote{}
	d.DB.First(&lastQuote)
	return lastQuote.Hash
}

func (d *DB) SetLastQuoteHash(hash string) {
	d.DB.Create(&LastQuote{Hash: hash})
}

func (d *DB) UpdateLastQuoteHash(hash string) {
	d.DB.First(&LastQuote{}).Update("hash", hash)
}
