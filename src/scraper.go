// Scraper logic
package scraper

import (
	"strings"

	"github.com/gocolly/colly"
)

type Page struct {
	Name       string
	URL        string
	VisitCheck bool
	Links      []*string
	Previous   *Page
	Depth      int
}

// Fungsi public untuk scrape
func PerformScrape(page *Page) {
	scrape(page)
}

func scrape(page *Page) {
	// Buat collector dengan domain khusus en.wikipedia.org
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML("title", func(e *colly.HTMLElement) {
		// Mengambil judul page
		page.Name = e.Text
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Mengambil URL page
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") {
			url := "https://en.wikipedia.org" + href
			page.Links = append(page.Links, &url)
		}
	})

	// Pergi ke URL page-nya biar bisa di-scrape
	c.Visit(page.URL)
}
