package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	"yada/internal"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	bot.Run(
		os.Getenv("YADA_TOKEN"),
		&internal.Bot{},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")
			return nil
		},
	)
}
