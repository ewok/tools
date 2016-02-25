// Package main provides ...
package packfreebook

import (
	"net/http"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func PackFreeBook() string {
	resp, err := http.Get("https://www.packtpub.com/packt/offers/free-learning")
	if err != nil {
		return "Error fetching url"
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "Error parsing"
	}

	book := cascadia.MustCompile(`div#deal-of-the-day div.dotd-main-book div.section-inner div.dotd-main-book-summary div.dotd-title h2 *`).MatchFirst(doc)
	header := strings.TrimSpace(book.Data)
	return header
}
