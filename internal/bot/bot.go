package bot

import (
	"math/rand"
	"time"

	"github.com/vitkarpenko/yada/internal/config"
	"github.com/vitkarpenko/yada/internal/services/emojis"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/services/muses"
	"github.com/vitkarpenko/yada/internal/services/reminders"
	"github.com/vitkarpenko/yada/internal/services/say"
	"github.com/vitkarpenko/yada/internal/utils"
	"github.com/vitkarpenko/yada/storages/sqlite"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Yada struct {
	Commands Commands
	Discord  *discordgo.Session
	Config   config.Config

	Queries *sqlite.Queries

	Images    *images.Service
	Emojis    *emojis.Service
	Muses     *muses.Service
	Say       *say.Service
	Reminders *reminders.Service
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't create discord session!")
	}

	yada := &Yada{
		Discord: discordSession,
		Config:  cfg,
	}

	yada.initDB()

	yada.initServices(discordSession, cfg)

	yada.setupIntents()

	return yada
}

func (y *Yada) initServices(discordSession *discordgo.Session, cfg config.Config) {
	y.Images = images.New(discordSession, cfg.ImagesChannelID, cfg.TenorAPIKey)
	y.Emojis = emojis.New(discordSession, cfg.GuildID)
	y.Muses = muses.New(discordSession, y.Queries, cfg)
	y.Reminders = reminders.New(discordSession, y.Queries, cfg)
	y.Say = say.New(y.Config.SoundsDataPath)
}

func (y *Yada) initDB() {
	db := sqlite.NewDB()
	y.Queries = sqlite.New(db)
}

func (y *Yada) Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := y.Discord.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't open websocket connection to discord!")
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
	y.Discord.AddHandler(func(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if c, ok := y.Commands[interaction.ApplicationCommandData().Name]; ok {
			c.Handler(s, interaction)
		}
	})

	// Add other handlers.
	y.Discord.AddHandler(y.AllMessagesHandler)
}
