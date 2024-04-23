package main

import (
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

func scrape(page *Page) {
	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		url := &href
		page.Links = append(page.Links, url)
	})

	//Placeholder
	c.Visit(page.URL)
}
