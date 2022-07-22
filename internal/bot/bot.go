package bot

import (
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/vitkarpenko/yada/internal/config"
	"github.com/vitkarpenko/yada/internal/services/emojis"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/utils"
)

type Yada struct {
	Commands Commands
	Discord  *discordgo.Session
	Config   config.Config

	Images *images.Service
	Emojis *emojis.Service
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	yada := &Yada{
		Discord: discordSession,
		Config:  cfg,
		Images:  images.New(discordSession, cfg.ImagesChannelID),
		Emojis:  emojis.New(discordSession, cfg.GuildID),
	}

	yada.setupIntents()

	return yada
}

func (y *Yada) Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := y.Discord.Open()
	if err != nil {
		log.Fatalln("Couldn't open websocket connection to discord!", err)
	}
	defer func(Discord *discordgo.Session) {
		_ = Discord.Close()
	}(y.Discord)

	y.setupInteractions()

	utils.WaitUntilInterrupted()
}

func (y *Yada) setupIntents() {
	y.Discord.Identify.Intents = discordgo.IntentsGuildMessages
}

func (y *Yada) setupInteractions() {
	y.setupCommands()
	y.setupHandlers()
}
func (y *Yada) setupCommands() {
	y.CleanupCommands()
	y.InitializeCommands()
}

func (y *Yada) setupHandlers() {
	// Add slash commands handlers.
	y.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := y.Commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	})

	// Add other handlers.
	y.Discord.AddHandler(y.AllMessagesHandler)
}
