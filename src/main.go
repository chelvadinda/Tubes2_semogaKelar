package main

import (
	"sync"
)

func BFS(startPage *Page) {
	var wg sync.WaitGroup
	queue := []*Page{startPage}
	visited := make(map[string]bool)
	depth := 0

	for len(queue) > 0 {
		for _, page := range queue {
			if visited[page.URL] {
				continue
			}

			visited[page.URL] = true

			wg.Add(1)
			go func(p *Page) {
				defer wg.Done()
				scrape(p)
			}(page)

			for _, link := range page.Links {
				childPage := &Page{
					Name:       *link,
					URL:        *link,
					VisitCheck: false,
					Previous:   page,
					Depth:      page.Depth + 1,
				}
				queue = append(queue, childPage)
			}
		}

		depth++
		queue = queue[len(queue):]

		wg.Wait()
	}
}

func IDS(startPage *Page) {
	var wg sync.WaitGroup
	depth := 0
	for {
		visited := make(map[string]bool)
		queue := []*Page{startPage}
		for len(queue) > 0 {
			page := queue[0]
			queue = queue[1:]

			if visited[page.URL] {
				continue
			}
			visited[page.URL] = true

			wg.Add(1)
			go func(p *Page) {
				defer wg.Done()
				scrape(p)
			}(page)

			for _, link := range page.Links {
				childPage := &Page{
					Name:       *link,
					URL:        *link,
					VisitCheck: false,
					Previous:   page,
					Depth:      page.Depth + 1,
				}
				queue = append(queue, childPage)
			}
		}
		wg.Wait()
		depth++
	}
}
