package bot

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const imagesPerReactionLimit = 5

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

			var files []*discordgo.File
			seen := make(map[string]bool)
			for i, word := range words {
				if image, ok := y.Images[word]; ok && !seen[word] {
					files = append(files, &discordgo.File{
						Name:        fmt.Sprintf("image_%d.gif", i),
						ContentType: "image/gif",
						Reader:      bytes.NewReader(image),
					})
					seen[word] = true
				}
			}

			if len(files) != 0 {
				files = limitFilesCount(files)
				_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Files: files,
				})
				if err != nil {
					log.Println("Couldn't send an image.", err)
				}
			}
		},
	)
}

func limitFilesCount(files []*discordgo.File) []*discordgo.File {
	if len(files) >= imagesPerReactionLimit {
		files = files[:imagesPerReactionLimit]
	}
	return files
}
