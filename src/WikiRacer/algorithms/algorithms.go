// Algorithms
package algorithms

import (
	"WikiRacer/scraper"
	"fmt"
	"strings"
	"time"
)

// Struct berupa hasil yang diberikan
type Result struct {
	SearchTime        time.Duration
	ArticlesChecked   int
	ArticlesTraversed int
	Route             []string
}

// Kondisi kedalaman yang paling dalam = 9 (cek QNA)
const maxDepth = 9

// Algoritma BFS
func BFS(startURL, targetURL string, algorithm string) Result {
	// Inisialisasi page awal
	source := &scraper.Page{
		URL:   startURL,
		Depth: 0,
	}

	var result Result
	startTime := time.Now()

	// Queue BFS
	queue := []*scraper.Page{source}

	// Map untuk cek halaman Wikipedia yang sudah dikunjungi
	visited := make(map[string]bool)
	visited[source.URL] = true

	// Algoritma BFS
	for len(queue) > 0 {
		currentPage := queue[0]
		queue = queue[1:]

		// Kondisi apabila target URL = URL yang sekarang
		if currentPage.URL == targetURL {
			result.Route = constructRoute(currentPage, algorithm)
			result.SearchTime = time.Since(startTime)
			return result
		}

		scraper.PerformScrape(currentPage)
		result.ArticlesChecked++

		for _, link := range currentPage.Links {
			childURL := *link
			if !visited[childURL] {
				childPage := &scraper.Page{
					Name:     getPageTitle(strings.TrimPrefix(childURL, "https://en.wikipedia.org/wiki/")),
					URL:      childURL,
					Previous: currentPage,
					Depth:    currentPage.Depth + 1,
				}
				queue = append(queue, childPage)
				visited[childURL] = true
			}
		}
	}

	// Return apabila tidak ditemukan
	result.SearchTime = time.Since(startTime)
	return result
}

// Fungsi membangun rute traversal artikel
func constructRoute(targetPage *scraper.Page, algorithm string) []string {
	var route []string
	current := targetPage

	for current != nil {
		pageTitle := current.Name
		if algorithm == "BFS" {
			pageTitle = strings.TrimSuffix(pageTitle, " - Wikipedia")
		}
		route = append(route, pageTitle)
		current = current.Previous
	}

	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}

	return route
}

// Algoritma IDS
func IDS(sourceURL, targetURL string) Result {
	source := &scraper.Page{
		URL:   sourceURL,
		Depth: 0,
	}

	var result Result
	startTime := time.Now()

	// Penerapan IDS
	for depth := 0; depth <= maxDepth; depth++ {
		fmt.Printf("Depth: %d\n", depth)
		result = DFS(source, targetURL, depth)
		if len(result.Route) > 0 {
			result.SearchTime = time.Since(startTime)
			return result
		}
	}

	result.SearchTime = time.Since(startTime)
	return result
}

// DFS sebagai kerangka IDS
func DFS(page *scraper.Page, targetURL string, depth int) Result {
	var result Result

	// Kondisi target URL berhasil ditemukan
	if page.URL == targetURL {
		fmt.Println("Target found:", page.URL)
		result.Route = append(result.Route, getPageTitle(page.Name))
		return result
	}
	if depth <= 0 {
		return result
	}

	// Proses scraping
	scraper.PerformScrape(page)
	result.ArticlesChecked++

	for _, link := range page.Links {
		childURL := *link
		childPage := &scraper.Page{
			Name:     getPageTitle(strings.TrimPrefix(childURL, "https://en.wikipedia.org/wiki/")),
			URL:      childURL,
			Previous: page,
			Depth:    page.Depth + 1,
		}

		childResult := DFS(childPage, targetURL, depth-1)
		result.ArticlesChecked += childResult.ArticlesChecked

		if len(childResult.Route) > 0 {
			result.Route = append(result.Route, getPageTitle(page.Name))
			result.Route = append(result.Route, childResult.Route...)
			return result
		}
	}

	return result
}

// Fungsi untuk mendapat nama page agar bisa dikeluarkan
func getPageTitle(pageName string) string {
	// Menghilangkan _ dari nama halaman
	pageName = strings.ReplaceAll(pageName, "_", " ")

	// Menghilangkan " - Wikipedia" dari nama halaman
	if idx := strings.Index(pageName, " - Wikipedia"); idx != -1 {
		return pageName[:idx]
	}
	return pageName
}
