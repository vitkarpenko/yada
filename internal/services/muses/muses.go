package muses

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/bwmarrin/discordgo"
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
		log.Println("Error while fetching muse rating: ", err)
	}

	var rating int
	switch {
	case err != sql.ErrNoRows:
		rating = int(savedMuseRating)
	case m.Author.ID == s.config.VitalyUserID:
		log.Println("A muse by Vitaly!")
		rating = normDistributedRating(10.8, 1)
	case m.Author.ID == s.config.LezhikUserID:
		log.Println("A muse by Lezhik!")
		rating = normDistributedRating(7.7, 1.4)
	case m.Author.ID == s.config.OlegUserID:
		log.Println("A muse by Oleg!")
		rating = normDistributedRating(5, 1.7)
	case m.Author.ID == s.config.VeraUserID:
		log.Println("A muse by Vera!")
		rating = 12
	default:
		log.Printf("A muse by %s!\n", m.Author.Username)
		rating = normDistributedRating(6, 1.8)
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
		log.Print("Error while creating a muse: ", err)
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
