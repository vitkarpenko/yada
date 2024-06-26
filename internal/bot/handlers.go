package bot

import (
	"github.com/vitkarpenko/yada/internal/tokens"
	"github.com/vitkarpenko/yada/internal/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const (
	randomImageChance = 0.01
	randomEmojiChance = 0.02
	randomGPTChance   = 0.05
)

func (y *Yada) AllMessagesHandler(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if ds.State.User != nil && m.Author.ID == ds.State.User.ID {
		return
	}

	y.handleImages(m)
	y.handleRandomEmoji(m)
	y.handleMuses(m)
	y.handleReminders(m)
}

func (y *Yada) handleRandomEmoji(m *discordgo.MessageCreate) {
	if utils.CheckChance(randomEmojiChance) {
		_, _ = y.Discord.ChannelMessageSendComplex(
			m.ChannelID,
			&discordgo.MessageSend{
				Content:   y.Emojis.Random(),
				Reference: y.refFromMessage(m),
			},
		)
	}
}

func (y *Yada) handleImages(m *discordgo.MessageCreate) {
	words := tokens.Tokenize(m.Content)

	if utils.CheckChance(randomImageChance) {
		randomImageURL, err := y.Images.RandomGifURL()
		if err != nil {
			log.Err(err).Msg("Error while fetching random gif")
			return
		}
		_, err = y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   randomImageURL,
			Reference: y.refFromMessage(m),
		})
		if err != nil {
			log.Error().Err(err).Msg("Couldn't send an image.")
		}
		return
	}

	files := y.Images.GetFilesToSend(words)
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

func (y *Yada) handleReminders(m *discordgo.MessageCreate) {
	y.Reminders.HandleMessage(m)
}

func (y *Yada) refFromMessage(m *discordgo.MessageCreate) *discordgo.MessageReference {
	return &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   y.Config.GuildID,
	}
}
