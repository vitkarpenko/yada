package quotes

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"yada/internal/config"
)

const (
	// 49385044 is my userID.
	quotesListURL = "https://www.goodreads.com/quotes/list/49385044?sort=date_added"
)

type Service struct {
	cfg    config.Goodreads
	client goodreadsClient
}

func NewService(cfg config.Goodreads) *Service {
	return &Service{cfg: cfg, client: newGoodreadsClient(cfg)}
}

type goodreadsClient struct {
	*http.Client
	sessionCookie string
}

func newGoodreadsClient(cfg config.Goodreads) goodreadsClient {
	client := goodreadsClient{
		Client:        http.DefaultClient,
		sessionCookie: cfg.SessionCookie,
	}
	client.getQuotes()

	return client
}

type Quote struct {
	text, author, book string
}

func (c goodreadsClient) getQuotes() []Quote {
	req := c.prepareQuotesListRequest()

	resp, err := c.Do(req)
	if err != nil {
		log.Println("Error while fetching quotes", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error while fetching quotes: status %d", resp.StatusCode)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error while parsing 'list quotes' response.")
	}

	authors := findAuthors(html)
	books := findBooks(html)
	quotes := findQuotes(html)
	if authors.Length() != books.Length() || books.Length() != quotes.Length() {
		log.Println("Got unmatched quotes, authors or book titles while fetching quotes.")
		return nil
	}

	quotesCount := quotes.Length()
	result := make([]Quote, quotesCount)
	for i := 0; i < quotesCount; i++ {
		quote := getQuote(quotes, authors, books, i)
		result[i] = quote
	}

	return result
}

func (c goodreadsClient) prepareQuotesListRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, quotesListURL, nil)

	q := req.URL.Query()
	q.Add("sort", "date_added")
	req.URL.RawQuery = q.Encode()

	return req
}
