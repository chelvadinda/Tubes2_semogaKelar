// Scraper logic
package scraper

import (
	"fmt"
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

func PerformScrape(page *Page) {
	scrape(page)
}

// Fungsi scrape
func scrape(page *Page) {
	// Buat collector dengan domain khusus en.wikipedia.org
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML("title", func(e *colly.HTMLElement) {
		// Mengambil judul page
		page.Name = e.Text
		fmt.Println("Title:", page.Name)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Mengambil URL page
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, "Main_Page") && !strings.Contains(href, "Wikipedia:") && !strings.Contains(href, "Portal:") && !strings.Contains(href, "Special:") && !strings.Contains(href, "Help:") && !strings.Contains(href, "Talk:") && !strings.Contains(href, "Category:") && !strings.Contains(href, "File:") && !strings.Contains(href, "Template:") && !strings.Contains(href, "Template_talk") {
			url := "https://en.wikipedia.org" + href
			page.Links = append(page.Links, &url)
			fmt.Println("Link:", url)
		}
	})

	// Pergi ke URL page-nya biar bisa di-scrape
	fmt.Println("Visiting:", page.URL)
	c.Visit(page.URL)
}

