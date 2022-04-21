package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/vitkarpenko/yada/internal/config"
)

const loadMessagesLimit = 100

type Yada struct {
	Commands Commands
	Discord  *discordgo.Session
	Images   map[string]Images
	Config   config.Config
	Emojis   []*discordgo.Emoji
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	yada := &Yada{
		Discord: discordSession,
		Images:  map[string]Images{},
		Config:  cfg,
	}

	yada.setupIntents()
	yada.getEmojis()

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

	y.startBackgroundTasks()
	y.setupInteractions()

	waitUntilInterrupted()
}

func (y *Yada) setupIntents() {
	y.Discord.Identify.Intents = discordgo.IntentsGuildMessages
}

func (y *Yada) getEmojis() {
	guildEmojis, err := y.Discord.GuildEmojis(y.Config.GuildID)
	if err != nil {
		log.Fatal("Couldn't get emojis!")
	}

	var availableEmojis []*discordgo.Emoji
	for _, e := range guildEmojis {
		if e.Available {
			availableEmojis = append(availableEmojis, e)
		}
	}

	y.Emojis = availableEmojis
}

func (y *Yada) setupInteractions() {
	y.setupCommands()
	y.setupHandlers()
}

func (y *Yada) startBackgroundTasks() {
	y.loadImagesInBackground()
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
	y.Discord.AddHandler(y.ReactWithImageHandler)
	y.Discord.AddHandler(y.RandomEmojiHandler)
}
