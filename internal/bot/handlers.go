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

			words := tokenize(m.Content)
			files := getFilesToSend(words, y.Images)

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

func getFilesToSend(words []string, images map[string]Image) []*discordgo.File {
	var files []*discordgo.File
	seenWords := make(map[string]bool)
	seenImages := make(map[string]bool)
	for _, word := range words {
		if seenWords[word] {
			continue
		}
		if image, ok := images[word]; ok {
			if seenImages[image.ID] {
				continue
			}
			files = append(files, &discordgo.File{
				Name:        fmt.Sprintf("image_%s.gif", image.ID),
				ContentType: "image/gif",
				Reader:      bytes.NewReader(image.Body),
			})
			seenImages[image.ID] = true
		}
		seenWords[word] = true
	}
	return files
}

func limitFilesCount(files []*discordgo.File) []*discordgo.File {
	if len(files) >= imagesPerReactionLimit {
		files = files[:imagesPerReactionLimit]
	}
	return files
}
