package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"yada/internal/bot"
	"yada/internal/config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not found")
	}

	var cfg config.Config
	envconfig.MustProcess("YADA", &cfg)

	yada := bot.NewYada(cfg)

	yada.Run()
}
