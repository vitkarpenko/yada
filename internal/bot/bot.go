package bot

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"yada/internal/config"
)

const loadMessagesLimit = 100

type Yada struct {
	Commands             Commands
	MessageReactHandlers []func(s *discordgo.Session, m *discordgo.MessageCreate)
	Discord              *discordgo.Session
	Images               map[string][]byte
	Config               config.Config
}

func NewYada(cfg config.Config) *Yada {
	discordSession, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalln("Couldn't create discord session!", err)
	}

	yada := &Yada{
		MessageReactHandlers: []func(s *discordgo.Session, m *discordgo.MessageCreate){},
		Discord:              discordSession,
		Images:               map[string][]byte{},
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

	y.loadImagesInBackground()

	y.setupCommands()
	y.setupHandlers()

	waitUntilInterrupted()
}

func (y *Yada) loadImagesInBackground() {
	y.processImages()
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for range ticker.C {
			y.processImages()
		}
	}()
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
	y.PrepareReactWithImageHandler()
	for _, handler := range y.MessageReactHandlers {
		y.Discord.AddHandler(handler)
	}
}

func (y *Yada) processImages() {
	var currentLastID string

	for {
		messages, err := y.Discord.ChannelMessages(
			y.Config.ImagesChannelID,
			loadMessagesLimit,
			"",
			currentLastID,
			"",
		)
		if err != nil {
			log.Fatalln("Could not load images from image channel!", err)
		}

		y.downloadImages(messages)

		if len(messages) < loadMessagesLimit {
			break
		}

		currentLastID = messages[len(messages)-1].ID
	}
}

func (y *Yada) downloadImages(messages []*discordgo.Message) {
	for _, message := range messages {
		attachments := message.Attachments
		if len(attachments) == 0 {
			continue
		}
		images := make([][]byte, len(attachments))
		for _, a := range attachments {
			images = append(images, readImageFromAttach(a))
		}
		triggerWords := strings.Split(message.Content, " ")
		y.setImagesTokens(triggerWords, images)
	}
}

func (y *Yada) setImagesTokens(triggerWords []string, images [][]byte) {
	for _, w := range triggerWords {
		for _, i := range images {
			y.Images[w] = i
		}
	}
}

func readImageFromAttach(a *discordgo.MessageAttachment) []byte {
	response, err := http.Get(a.URL)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	image, _ := io.ReadAll(response.Body)
	return image
}
