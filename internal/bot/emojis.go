package bot

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

const randomEmojiChance = 0.05

var randomEmojis = []string{
	"<:peka:710107629761855619>",
	"<:uuu:884497957108457532>",
	"<:epeka:912368988615495710>",
	"<:smlpeka:882817339614199889>",
	"<:yao:882718429667295304>",
	"<:tears:405771172903649280>",
	"<:thinking2:882730590590361672>",
	"<:gusta:882722295783780383>",
	"<:why:882720049658470400>",
}

func (y *Yada) RandomEmojiHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if checkChance(randomEmojiChance) {
		_, _ = y.Discord.ChannelMessageSendComplex(
			m.ChannelID,
			&discordgo.MessageSend{
				Content: getRandomEmoji(),
				Reference: &discordgo.MessageReference{
					MessageID: m.Message.ID,
					ChannelID: m.ChannelID,
					GuildID:   y.Config.GuildID,
				},
			},
		)
	}
}

func getRandomEmoji() string {
	return randomEmojis[rand.Intn(len(randomEmojis))]
}
