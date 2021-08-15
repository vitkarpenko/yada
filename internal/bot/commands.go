package bot

import "github.com/bwmarrin/discordgo"

type Command struct {
	AppCommand discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Commands maps command names to Command instances.
type Commands map[string]Command

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
}
