package sqlite

import (
	"database/sql"
	"io/fs"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func SetEmbedMigrations(embed fs.FS) {
	goose.SetBaseFS(embed)
}

func NewDB() *sql.DB {

	db, err := sql.Open("sqlite3", "data/sqlite.db")
	if err != nil {
		log.Fatal("Error while creating db: ", err)
	}

	migrate(db)

	return db
}

func migrate(db *sql.DB) {
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
