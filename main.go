package main

import (
	"embed"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"github.com/vitkarpenko/yada/internal/bot"
	"github.com/vitkarpenko/yada/internal/config"
	"github.com/vitkarpenko/yada/storages/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msg(".env not found")
	}

	var cfg config.Config
	envconfig.MustProcess("YADA", &cfg)

	sqlite.SetEmbedMigrations(embedMigrations)

	yada := bot.NewYada(cfg)

	yada.Run()
}
