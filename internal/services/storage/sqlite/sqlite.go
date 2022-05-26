package sqlite

import (
	"errors"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const dbPath = "data/sqlite.db"

type DB struct {
	*gorm.DB
}

func New() *DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&Swear{}); err != nil {
		log.Fatal(err)
	}

	return &DB{db}
}

func (db *DB) UploadSwears(swears []string) error {
	swearRows := make([]Swear, len(swears))
	for i, s := range swears {
		swearRows[i] = Swear{Word: s}
	}

	result := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(swearRows, 1000)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) ShouldFillSwears() bool {
	var count int64
	db.Model(&Swear{}).Count(&count)

	return count == 0
}

func (db *DB) HasSwear(words []string) (string, error) {
	swear := &Swear{}
	err := db.Where("word IN ?", words).First(swear).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}

	return swear.Word, nil
}
