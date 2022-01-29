package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (y *Yada) ChoiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
