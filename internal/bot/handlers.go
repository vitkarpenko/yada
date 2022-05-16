package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/tokens"
	"github.com/vitkarpenko/yada/internal/utils"
)

const (
	randomImageChance = 0.005
	randomEmojiChance = 0.025
)

func (y *Yada) ReactWithImageHandler(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == ds.State.User.ID {
		return
	}

	words := tokens.Tokenize(m.Content)

	var files []*discordgo.File
	if utils.CheckChance(randomImageChance) {
		files = []*discordgo.File{
			images.DiscordFileFromImage(y.Images.Random(), uuid.New().String()),
		}
	} else {
		files = y.Images.GetFilesToSend(words)
		fmt.Println(files)
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

func (y *Yada) RandomEmojiHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if utils.CheckChance(randomImageChance) {
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
