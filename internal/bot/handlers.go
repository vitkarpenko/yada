package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Choice(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func ReactWithImage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	_ = tokenize(message)
}
