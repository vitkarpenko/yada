package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"

	"yada/internal/bot"
	"yada/internal/config"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config.Config
	envconfig.MustProcess("YADA", &cfg)

	yada := bot.NewYada(cfg)

	yada.Run()
}
