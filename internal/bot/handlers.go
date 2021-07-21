package bot

import (
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

			var imageEmbed *discordgo.MessageEmbedImage
			for _, word := range words {
				if image, ok := y.Images[word]; ok {
					imageEmbed = &discordgo.MessageEmbedImage{
						URL:      image.URL,
						ProxyURL: image.ProxyURL,
						Width:    image.Width,
						Height:   image.Height,
					}
					break
				}
			}

			if imageEmbed != nil {
				_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Embed: &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{
							IconURL: discordgo.EndpointUserAvatar(m.Author.ID, m.Author.Avatar),
							Name:    m.Author.Username,
						},
						Type:  discordgo.EmbedTypeImage,
						Image: imageEmbed,
					},
				})
				if err != nil {
					log.Printf("Couldn't send an image.", err)
				}
			}
		},
	)
}
