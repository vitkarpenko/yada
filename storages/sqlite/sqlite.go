package sqlite

import (
	"database/sql"
	"io/fs"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func SetEmbedMigrations(embed fs.FS) {
	goose.SetBaseFS(embed)
}

func NewDB() *sql.DB {

	db, err := sql.Open("sqlite3", "data/sqlite.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating db")
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
