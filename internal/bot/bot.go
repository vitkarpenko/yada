package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"yada/internal/services/quotes"
	"yada/internal/storage/postgres"

	"yada/internal/config"
)

const loadMessagesLimit = 100

type Yada struct {
	Commands  Commands
	Discord   *discordgo.Session
	DB        *postgres.DB
	Images    map[string]Images
	Reminders []postgres.Reminder
	Config    config.Config
	Quotes    *quotes.Service
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
		Discord: discordSession,
		DB:      db,
		Images:  map[string]Images{},
		Config:  cfg,
		Quotes:  quotes.NewService(cfg.Quotes, discordSession, db),
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

	y.initialize()
	y.startBackgroundTasks()
	y.setupInteractions()

	waitUntilInterrupted()
}

func (y *Yada) setupInteractions() {
	y.setupCommands()
	y.setupHandlers()
}

func (y *Yada) startBackgroundTasks() {
	y.checkRemindersInBackground()
	y.loadImagesInBackground()
	y.Quotes.CheckQuotesInBackground()
}

func (y *Yada) initialize() {
	y.loadReminders()
}

func (y *Yada) setupIntents() {
	y.Discord.Identify.Intents = discordgo.IntentsGuildMessages
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
	y.Discord.AddHandler(y.SetReminderHandler)
}
