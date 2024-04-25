// Algorithms
package algorithms

import (
	"WikiRacer/scraper"
	"sync"
	"time"
)

// Struct berupa hasil yang diberikan
type Result struct {
	SearchTime        time.Duration
	ArticlesChecked   int
	ArticlesTraversed int
	Route             []string
}

// Algoritma BFS
func BFS(startPage *scraper.Page, targetURL string, result *Result) {
	// Mulai waktu
	start := time.Now()

	var wg sync.WaitGroup               // Memastikan semua Goroutine berjalan dengan teratur
	queue := []*scraper.Page{startPage} // Queue untuk menjalankan BFS
	visited := make(map[string]bool)    // Map untuk menghindari page yang sudah dilalui
	depth := 0                          // Inisialisasi depth untuk kondisi terminasi
	linksChecked := 0                   // Bagian dari hasil
	articlesTraversed := 0              // Bagian dari hasil

	// Loop sepanjang queue dan berhenti kalau udah kedalaman = 9 (Cek QNA untuk referensi)
	// Kondisi terminasi agar tidak terlalu lama melakukan proses-nya
	for len(queue) > 0 && depth < 9 {
		for _, page := range queue {
			articlesTraversed++

			if visited[page.URL] {
				continue
			}

			visited[page.URL] = true

			result.Route = getPath(page)

			// Cek URL saat ini dengan targetURL
			if page.URL == targetURL {
				result.SearchTime = time.Since(start)
				result.ArticlesChecked = linksChecked
				result.ArticlesTraversed = articlesTraversed
				result.Route = getPath(page)
				return
			}

			// wg.Add(1) untuk menambah proses Goroutine yang dilakukan
			wg.Add(1)
			go func(p *scraper.Page) {
				defer wg.Done()
				scraper.PerformScrape(p)
			}(page)

			linksChecked += len(page.Links)

			// Inisialisasi untuk setiap hyperlink di suatu page untuk dikunjungi
			for _, link := range page.Links {
				childPage := &scraper.Page{
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

		// Blok fungsi utama hingga semua proses selesai dilakukan
		wg.Wait()
	}

	// Mengirim hasil kosong apabila sudah terminasi dan tidak menemukan targetURL yang dicari
	result.SearchTime = time.Since(start)
	result.ArticlesChecked = linksChecked
	result.ArticlesTraversed = articlesTraversed
	result.Route = nil
}

// Algoritma IDS
func IDS(startPage *scraper.Page, targetURL string, result *Result) {
	// Mulai waktu
	start := time.Now()

	var wg sync.WaitGroup  // Memastikan semua Goroutine berjalan dengan teratur
	depth := 0             // Inisialisasi kedalaman
	linksChecked := 0      // Bagian dari hasil
	articlesTraversed := 0 // Bagian dari hasil

	// Loop hingga kedalaman = 9 (Cek QNA untuk referensi)
	// Kondisi terminasi agar tidak terlalu lama melakukan proses-nya
	for depth < 9 {
		visited := make(map[string]bool)
		queue := []*scraper.Page{startPage}
		for len(queue) > 0 {
			page := queue[0]
			queue = queue[1:]

			articlesTraversed++

			if visited[page.URL] {
				continue
			}
			visited[page.URL] = true

			result.Route = getPath(page)

			// Cek URL saat ini dengan targetURL
			if page.URL == targetURL {
				result.SearchTime = time.Since(start)
				result.ArticlesChecked = linksChecked
				result.ArticlesTraversed = articlesTraversed
				result.Route = getPath(page)
				return
			}

			// wg.Add(1) untuk menambah proses Goroutine yang dilakukan
			wg.Add(1)
			go func(p *scraper.Page) {
				defer wg.Done()
				scraper.PerformScrape(p)
			}(page)

			linksChecked += len(page.Links)

			for _, link := range page.Links {
				childPage := &scraper.Page{
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

	// Mengirim hasil kosong apabila sudah terminasi dan tidak menemukan targetURL yang dicari
	result.SearchTime = time.Since(start)
	result.ArticlesChecked = linksChecked
	result.ArticlesTraversed = articlesTraversed
	result.Route = nil
}

// Fungsi untuk mengirim path atau jalur dari page yang dilalui
func getPath(page *scraper.Page) []string {
	var path []string
	current := page

	// Loop untuk append page dengan yang lebih akhir duluan
	for current != nil {
		path = append(path, current.Name)
		current = current.Previous
	}
	// Membalikkan unsur path agar bisa menunjukkan page dari awal ke akhir
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
