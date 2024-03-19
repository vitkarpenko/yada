package emojis

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	emojis  []*discordgo.Emoji
	discord *discordgo.Session
	guildID string
}

func New(discord *discordgo.Session, guildID string) *Service {
	service := &Service{
		discord: discord,
		guildID: guildID,
	}
	service.getEmojis()
	return service
}

func (s *Service) Random() string {
	if len(s.emojis) == 0 {
		return ""
	}

	emoji := s.emojis[rand.Intn(len(s.emojis))]
	return fmt.Sprintf("<:%s:%s>", emoji.Name, emoji.ID)
}

func (s *Service) getEmojis() {
	guildEmojis, err := s.discord.GuildEmojis(s.guildID)
	if err != nil {
		log.Fatal("Couldn't get emojis!")
	}

	var availableEmojis []*discordgo.Emoji
	for _, e := range guildEmojis {
		if e.Available || !e.Animated {
			availableEmojis = append(availableEmojis, e)
		}
	}

	s.emojis = availableEmojis
}
