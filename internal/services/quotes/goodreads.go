package quotes

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type goodreadsClient struct {
	*http.Client
}

func newGoodreadsClient() *goodreadsClient {
	client := goodreadsClient{
		Client: http.DefaultClient,
	}

	return &client
}

func (gc goodreadsClient) getQuotes() []Quote {
	req := gc.prepareQuotesListRequest()

	resp, err := gc.Do(req)
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
		quote.hash = md5Hash(quote.text)
		result[i] = quote
	}

	return result
}

func (gc goodreadsClient) prepareQuotesListRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, quotesListURL, nil)

	q := req.URL.Query()
	q.Add("sort", "date_added")
	req.URL.RawQuery = q.Encode()

	return req
}
