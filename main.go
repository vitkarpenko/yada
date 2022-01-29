package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/joho/godotenv"

	"yada/internal/bot"
	"yada/internal/config"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not found")
	}

	var cfg config.Config
	envconfig.MustProcess("YADA", &cfg)

	yada := bot.NewYada(cfg)

	yada.Run()
}
