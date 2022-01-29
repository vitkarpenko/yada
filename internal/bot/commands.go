package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
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
		log.Fatal("Cannot fetch commands.", err)
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
				Description: "Выбираю для тебя случайный элемент из списка, неуверенный кожаный мешок.",
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
	}

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
