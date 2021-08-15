package bot

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
	"yada/internal/storage/postgres"

	"yada/internal/config"
)

const loadMessagesLimit = 100

type Yada struct {
	Commands             Commands
	MessageReactHandlers []func(s *discordgo.Session, m *discordgo.MessageCreate)
	Discord              *discordgo.Session
	DB                   *gorm.DB
	Images               map[string]Image
	Reminders            []postgres.Reminder
	Config               config.Config
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalln("Couldn't connect to database!", err)
	}

	yada := &Yada{
		MessageReactHandlers: []func(s *discordgo.Session, m *discordgo.MessageCreate){},
		Discord:              discordSession,
		DB:                   db,
		Images:               map[string]Image{},
		Config:               cfg,
	}

	yada.setupIntents()

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

	y.loadReminders()
	y.checkRemindersInBackground()

	y.loadImagesInBackground()

	y.setupCommands()
	y.setupHandlers()

	waitUntilInterrupted()
}

func (y *Yada) setupIntents() {
	y.Discord.Identify.Intents = discordgo.IntentsGuildMessages
}

func (y *Yada) setupCommands() {
	y.InitializeCommands()

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

func (y *Yada) setupHandlers() {
	// Add slash commands handlers.
	y.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := y.Commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	})

	// Add other handlers.
	y.Discord.AddHandler(y.ReactWithImageHandler)
	y.Discord.AddHandler(y.SetReminderHandler)
}
