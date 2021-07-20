package commands

import "github.com/bwmarrin/discordgo"

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "kek",
			Description: "Random kek.",
		},
	}

	Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"choice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Slash command",
				},
			})
		},
	}
)
