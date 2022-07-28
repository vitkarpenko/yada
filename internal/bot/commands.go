package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Command struct {
	AppCommand discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Commands maps command names to Command instances.
type Commands map[string]Command

// CleanupCommands deletes all existing commands. Kinda overkill, but who cares. :)
func (y *Yada) CleanupCommands() {
	commands, err := y.Discord.ApplicationCommands(y.Config.AppID, y.Config.GuildID)
	if err != nil {
		log.Fatal().Msg("Cannot fetch commands!")
	}
	for _, c := range commands {
		_ = y.Discord.ApplicationCommandDelete(y.Config.AppID, y.Config.GuildID, c.ID)
	}
}

func (y *Yada) InitializeCommands() {
	y.Commands = Commands{
		"choice": Command{
			AppCommand: discordgo.ApplicationCommand{
				Name:        "choice",
				Description: "Выбираю для тебя случайный элемент из списка.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "варианты",
						Description: "Список вариантов, разделённых запятой.",
						Required:    true,
					},
				},
			},
			Handler: y.ChoiceHandler,
		},
		"gachi": Command{
			AppCommand: discordgo.ApplicationCommand{
				Name:        "gachi",
				Description: "Oh shit, I'm sorry!",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "звук",
						Description:  "Что проорать?",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			Handler: y.Gachi.Handler,
		},
	}

	for _, c := range y.Commands {
		appCommand := &c.AppCommand
		_, err := y.Discord.ApplicationCommandCreate(
			y.Discord.State.User.ID,
			y.Config.GuildID,
			appCommand,
		)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create '%v' command: %v", appCommand.Name)
		}
	}
}

func (y *Yada) ChoiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.ApplicationCommandData().Options[0].StringValue()
	words := strings.Split(message, ",")
	randIndex := rand.Intn(len(words))

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Выбирала из списка `%s` и выбрала: **%s**.",
				message,
				strings.TrimSpace(words[randIndex]),
			),
		},
	})
}
