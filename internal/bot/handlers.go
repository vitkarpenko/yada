package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/tokens"
	"github.com/vitkarpenko/yada/internal/utils"
)

const (
	randomImageChance = 0.01
	randomEmojiChance = 0.02
)

func (y *Yada) AllMessagesHandler(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if ds.State.User != nil && m.Author.ID == ds.State.User.ID {
		return
	}

	isSwear := y.handleSwears(m)
	if isSwear {
		return
	}
	y.handleImages(m)
	y.handleRandomEmoji(m)
}

func (y *Yada) handleRandomEmoji(m *discordgo.MessageCreate) {
	if utils.CheckChance(randomEmojiChance) {
		_, _ = y.Discord.ChannelMessageSendComplex(
			m.ChannelID,
			&discordgo.MessageSend{
				Content: y.Emojis.Random(),
				Reference: &discordgo.MessageReference{
					MessageID: m.Message.ID,
					ChannelID: m.ChannelID,
					GuildID:   y.Config.GuildID,
				},
			},
		)
	}
}

func (y *Yada) handleImages(m *discordgo.MessageCreate) {
	words := tokens.Tokenize(m.Content)

	var files []*discordgo.File
	if utils.CheckChance(randomImageChance) {
		files = []*discordgo.File{
			images.DiscordFileFromImage(y.Images.Random(), uuid.New().String()),
		}
	} else {
		files = y.Images.GetFilesToSend(words)
	}

	if len(files) != 0 {
		_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: files,
		})
		if err != nil {
			log.Println("Couldn't send an image.", err)
		}
	}
}

func (y *Yada) handleSwears(m *discordgo.MessageCreate) (isSwear bool) {
	words := tokens.Tokenize(m.Content)

	for _, w := range words {
		if y.Swears.IsSwear(w) {
			capitalizedWord := strings.Title(strings.ToLower(w))

			_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Files: []*discordgo.File{
					images.DiscordFileFromImage(y.Swears.PunishmentImage(), uuid.New().String()),
				},
				Content:   fmt.Sprintf("%s? %s", capitalizedWord, y.Swears.PunishmentPhrase()),
				Reference: m.Reference(),
			})
			if err != nil {
				log.Println("Couldn't send swear punishment.", err)
			}

			return true
		}
	}

	return
}
