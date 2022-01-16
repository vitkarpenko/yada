package quotes

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getQuote(
	quotes *goquery.Selection,
	authors *goquery.Selection,
	books *goquery.Selection,
	i int,
) Quote {
	return Quote{
		text:   getQuoteText(quotes, i),
		author: getAuthor(authors, i),
		book:   getBook(books, i),
	}
}

func findQuotes(html *goquery.Document) *goquery.Selection {
	return html.Find(".quoteText")
}

func findBooks(html *goquery.Document) *goquery.Selection {
	return html.Find("span > .authorOrTitle")
}

func findAuthors(html *goquery.Document) *goquery.Selection {
	return html.Find("div > .authorOrTitle")
}

func getQuoteText(quotes *goquery.Selection, i int) string {
	return strings.TrimSpace(quotes.Get(i).FirstChild.Data)
}

func getBook(books *goquery.Selection, index int) string {
	return strings.Trim(books.Get(index).FirstChild.Data, " ,\n")
}

func getAuthor(authors *goquery.Selection, i int) string {
	return strings.Trim(authors.Get(i).FirstChild.Data, " ,\n")
}
