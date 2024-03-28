package gpt

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const discordMaximumMsgLength = 2000

type Service struct {
	discord *discordgo.Session
	client  client
}

func New(session *discordgo.Session, token, catalogID string) *Service {
	return &Service{
		discord: session,
		client:  newClient(token, catalogID),
	}
}

func (s *Service) HandleMessage(m *discordgo.MessageCreate) {
	var resp string

	content := filterMentions(m.Content)

	reply, err := s.client.reply(filterMentions(content))
	if err != nil {
		resp = fmt.Sprintf("Ошибка при запросе в Yandex GPT: %v", err)
	} else {
		resp = reply
	}

	runes := []rune(resp)
	start := 0
	for {
		end := start + discordMaximumMsgLength
		if end > len(runes) {
			end = len(runes)
		}

		_, _ = s.discord.ChannelMessageSend(
			m.ChannelID,
			string(runes[start:end]),
		)

		if end == len(runes) {
			break
		}

		start = end
	}
}

func filterMentions(s string) string {
	tokens := strings.Split(s, " ")

	filtered := slices.DeleteFunc(tokens, func(t string) bool {
		return strings.Contains(t, "@")
	})

	return strings.Join(filtered, " ")
}
