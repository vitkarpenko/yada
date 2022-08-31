package muses

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/vitkarpenko/yada/internal/config"
	"github.com/vitkarpenko/yada/storages/sqlite"
)

type Service struct {
	discord *discordgo.Session
	queries *sqlite.Queries

	config config.Config
}

func New(discord *discordgo.Session, queries *sqlite.Queries, config config.Config) *Service {
	return &Service{discord: discord, queries: queries, config: config}
}

var (
	MehMuseEmojis     = []string{"(＃＞＜)", "<(￣ ﹌ ￣)>", "(＞﹏＜)", "ヾ( ￣O￣)ツ"}
	FineMuseEmojis    = []string{"(￣_￣)・・・", "(•ิ_•ิ)?", "┐(￣ヘ￣;)┌", "ლ(ಠ_ಠ ლ)"}
	AwesomeMuseEmojis = []string{"^ - ^", "(♡-_-♡)", "ヽ(♡‿♡)ノ", "♡( ◡‿◡ )"}
	WaifuMuseEmojis   = []string{"o(≧▽≦)o", "(❤ω❤)", "♡＼(￣▽￣)／♡", "(*♡∀♡)"}
)

func (s *Service) HandleMessage(m *discordgo.MessageCreate) {
	if len(m.Attachments) > 1 {
		_, _ = s.discord.ChannelMessageSend(
			m.ChannelID,
			"Не могу выразить своё экспертное мнение, слишком много муз! 🤌🏻",
		)
	}

	hash := imageHash(m)
	savedMuseRating, err := s.queries.GetMuseRating(context.Background(), hash)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("Error while fetching muse rating!")
	}

	var rating int
	if err != sql.ErrNoRows {
		rating = int(savedMuseRating)
	} else {
		rating = normDistributedRating(7.7, 1.8)
	}

	var (
		punctuation string
		emojis      []string
	)
	switch {
	case 0 <= rating && rating <= 3:
		emojis = MehMuseEmojis
		punctuation = "..."
	case 4 <= rating && rating <= 7:
		emojis = FineMuseEmojis
		punctuation = "."
	case 8 <= rating && rating <= 10:
		emojis = AwesomeMuseEmojis
		punctuation = "!"
	case 11 <= rating && rating <= 12:
		emojis = WaifuMuseEmojis
		punctuation = "!!!"
	}

	emoji := emojis[rand.Intn(len(emojis))]
	message := fmt.Sprintf("%d/10%s %s", rating, punctuation, emoji)

	_, _ = s.discord.ChannelMessageSendComplex(
		m.ChannelID,
		&discordgo.MessageSend{
			Content: message,
			Reference: &discordgo.MessageReference{
				MessageID: m.Message.ID,
				ChannelID: m.ChannelID,
				GuildID:   s.config.GuildID,
			},
		},
	)

	err = s.queries.CreateMuse(
		context.Background(),
		sqlite.CreateMuseParams{
			Hash:   hash,
			Rating: int64(rating),
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating a muse: ")
	}
}

func imageHash(m *discordgo.MessageCreate) string {
	image := readImageBodyFromAttach(m.Attachments[0])
	hash := md5.Sum(image)
	return hex.EncodeToString(hash[:])
}

func readImageBodyFromAttach(a *discordgo.MessageAttachment) []byte {
	response, err := http.Get(a.URL)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	imageBody, _ := io.ReadAll(response.Body)

	return imageBody
}

func normDistributedRating(offset, stdev float64) int {
	rating := rand.NormFloat64()*stdev + offset
	switch {
	case rating < 0:
		return 0
	case rating > 12:
		return 12
	default:
		return int(rating)
	}
}
