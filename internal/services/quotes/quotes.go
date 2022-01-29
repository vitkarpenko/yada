package quotes

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"yada/internal/config"
	"yada/internal/storage/postgres"
)

const (
	// 49385044 is my userID.
	quotesListURL = "https://www.goodreads.com/quotes/list/49385044?sort=date_added"
)

type Quote struct {
	text, author, book string
	hash               string
}

type Service struct {
	cfg     config.Quotes
	discord *discordgo.Session
	db      *postgres.DB
	client  *goodreadsClient
}

func NewService(cfg config.Quotes, discord *discordgo.Session, db *postgres.DB) *Service {
	return &Service{cfg: cfg, discord: discord, db: db, client: newGoodreadsClient()}
}

func (s *Service) CheckQuotesInBackground() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			s.checkQuotes()
		}
	}()
}

func (s *Service) checkQuotes() {
	quotes := s.getNewQuotes()
	// Post old quotes first.
	quotes = reverseQuotes(quotes)

	for _, q := range quotes {
		for {
			err := s.postQuote(q)
			if err != nil {
				log.Println("Error while sending quote to discord", err)
				return
			} else {
				break
			}
		}
	}
}

func (s *Service) getNewQuotes() (result []Quote) {
	lastQuoteHash := s.db.GetLastQuoteHash()
	quotes := s.client.getQuotes()

	if lastQuoteHash == "" && len(quotes) >= 1 {
		s.db.SetLastQuoteHash(quotes[0].hash)
		return quotes
	}

	for _, q := range quotes {
		if q.hash == lastQuoteHash {
			break
		}
		result = append(result, q)
	}

	if len(result) >= 1 {
		s.db.UpdateLastQuoteHash(result[0].hash)
	}

	return
}

func (s *Service) postQuote(q Quote) error {
	message := formatQuoteMessage(q)
	_, err := s.discord.ChannelMessageSend(s.cfg.QuotesChannelID, message)
	return err
}

func formatQuoteMessage(q Quote) string {
	message := strings.Join(
		[]string{
			fmt.Sprintf("> %s", q.text),
			fmt.Sprintf("*— %s, %s*", q.author, q.book),
		},
		"\n",
	)
	return message
}

func md5Hash(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func reverseQuotes(quotes []Quote) []Quote {
	for i, j := 0, len(quotes)-1; i < j; i, j = i+1, j-1 {
		quotes[i], quotes[j] = quotes[j], quotes[i]
	}
	return quotes
}
