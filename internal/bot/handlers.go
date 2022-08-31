package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/tokens"
	"github.com/vitkarpenko/yada/internal/utils"
)

const (
	randomImageChance = 0.1
	randomEmojiChance = 0.04
)

func (y *Yada) AllMessagesHandler(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if ds.State.User != nil && m.Author.ID == ds.State.User.ID {
		return
	}

	y.handleImages(m)
	y.handleRandomEmoji(m)
	y.handleMuses(m)
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
		randomImage, err := y.Images.Random()
		if err != nil {
			log.Err(err).Msg("Error while fetching random gif")
			return
		}
		files = []*discordgo.File{
			images.DiscordFileFromImage(randomImage, "video/mp4"),
		}
	} else {
		files = y.Images.GetFilesToSend(words)
	}

	if len(files) != 0 {
		_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: files,
		})
		if err != nil {
			log.Error().Err(err).Msg("Couldn't send an image.")
		}
	}
}

func (y *Yada) handleMuses(m *discordgo.MessageCreate) {
	if m.ChannelID != y.Config.MusesChannelID || len(m.Attachments) == 0 {
		return
	}

	y.Muses.HandleMessage(m)
}
