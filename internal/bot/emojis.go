package bot

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

const randomEmojiChance = 0.04

func (y *Yada) RandomEmojiHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if checkChance(randomEmojiChance) {
		_, _ = y.Discord.ChannelMessageSendComplex(
			m.ChannelID,
			&discordgo.MessageSend{
				Content: y.getRandomEmoji(),
				Reference: &discordgo.MessageReference{
					MessageID: m.Message.ID,
					ChannelID: m.ChannelID,
					GuildID:   y.Config.GuildID,
				},
			},
		)
	}
}

func (y *Yada) getRandomEmoji() string {
	emoji := y.Emojis[rand.Intn(len(y.Emojis))]
	return fmt.Sprintf("<:%s:%s>", emoji.Name, emoji.ID)
}
