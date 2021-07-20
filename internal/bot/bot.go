package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"

	"yada/internal/config"
)

type Yada struct {
	Commands Commands
	Discord  *discordgo.Session
	Images   map[string]*discordgo.MessageAttachment
	Config   config.Config
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	yada := &Yada{
		Commands: InitializeCommands(),
		Discord:  discordSession,
		Images:   map[string]*discordgo.MessageAttachment{},
		Config:   cfg,
	}

	yada.setupIntents()
	yada.loadImages()

	return yada
}

func (y *Yada) Run() {
	err := y.Discord.Open()
	if err != nil {
		log.Fatalln("Couldn't open websocket connection to discord!", err)
	}
	defer func(Discord *discordgo.Session) {
		_ = Discord.Close()
	}(y.Discord)

	y.setupCommands()
	y.setupHandlers()

	waitUntilInterrupted()
}

func (y *Yada) setupIntents() {
	y.Discord.Identify.Intents = discordgo.IntentsGuildMessages
}

func (y *Yada) setupCommands() {
	for _, c := range y.Commands {
		appCommand := &c.AppCommand
		_, err := y.Discord.ApplicationCommandCreate(
			y.Discord.State.User.ID,
			y.Config.GuildID,
			appCommand,
		)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", appCommand.Name, err)
		}
	}
}

func (y *Yada) setupHandlers() func() {
	return y.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := y.Commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	})
}

func (y *Yada) loadImages() {
	messages, err := y.Discord.ChannelMessages(y.Config.ImagesChannelID, 0, "", "", "")
	if err != nil {
		log.Fatalln("Could not load images from image channel!", err)
	}

	for _, message := range messages {
		triggerWords := strings.Split(message.Content, " ")
		for _, w := range triggerWords {
			y.Images[w] = message.Attachments[0]
		}
	}
}
