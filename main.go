package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"yada/internal/bot"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("YADA_TOKEN")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		log.Fatalln("Couldn't open websocket connection to discord!", err)
	}
	defer discord.Close()

	commands := bot.InitializeCommands()
	setupCommands(discord, commands)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func setupCommands(discord *discordgo.Session, commands bot.Commands) {
	createApplicationCommands(discord, commands)
	setupHandlers(discord, commands)
}

func createApplicationCommands(discord *discordgo.Session, commands bot.Commands) {
	for _, c := range commands {
		slashCommand := &c.AppCommand
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, os.Getenv("GUILD_ID"), slashCommand)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", slashCommand.Name, err)
		}
	}
}

func setupHandlers(discord *discordgo.Session, commands bot.Commands) func() {
	return discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	})
}
