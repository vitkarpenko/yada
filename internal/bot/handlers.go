package bot

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (y *Yada) PrepareChoiceHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		message := i.ApplicationCommandData().Options[0].StringValue()
		words := strings.Split(message, ",")
		randIndex := rand.Intn(len(words))

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Выбирал из списка `%s` и выбрал: **%s**.",
					message,
					strings.TrimSpace(words[randIndex]),
				),
			},
		})
	}
}

func (y *Yada) PrepareReactWithImageHandler() {
	y.MessageReactHandlers = append(
		y.MessageReactHandlers,
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			// Ignore all messages created by the bot itself.
			if m.Author.ID == s.State.User.ID {
				return
			}

			message := m.Content
			words := tokenize(message)

			var file *discordgo.File
			for _, word := range words {
				if image, ok := y.Images[word]; ok {
					file = &discordgo.File{
						Name:        "image.gif",
						ContentType: "image/gif",
						Reader:      bytes.NewReader(image),
					}
					break
				}
			}

			if file != nil {
				_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					File: file,
				})
				if err != nil {
					log.Println("Couldn't send an image.", err)
				}
			}
		},
	)
}
