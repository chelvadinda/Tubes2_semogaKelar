package main

import (
	"WikiRacer/algorithms"
	"WikiRacer/scraper"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve index.html
	http.ServeFile(w, r, "index.html")
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	// Buat file server handler
	fs := http.FileServer(http.Dir("./"))

	// Serve CSS file
	fs.ServeHTTP(w, r)
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	startURLPart, targetURLPart := extractFormData(r)

	// Buat start dan target artikel serta ambil algoritma
	startArticle := scraper.Page{URL: "https://en.wikipedia.org/wiki/" + startURLPart}
	targetArticle := scraper.Page{URL: "https://en.wikipedia.org/wiki/" + targetURLPart}

	algorithm := r.FormValue("algorithm")

	// Panggil fungsi-fungsi
	var result algorithms.Result
	switch algorithm {
	case "bfs":
		result = algorithms.BFS(startArticle.URL, targetArticle.URL, algorithm, &startArticle)
		result.ArticlesTraversed = len(result.Route)
	case "ids":
		result = algorithms.IDS(startArticle.URL, targetArticle.URL)
		result.ArticlesTraversed = len(result.Route)
	default:
		http.Error(w, "Invalid algorithm selection", http.StatusBadRequest)
		return
	}

	fmt.Println("Search Time:", result.SearchTime)
	fmt.Println("Articles Checked:", result.ArticlesChecked)
	fmt.Println("Articles Traversed:", len(result.Route))
	fmt.Println("Route:", strings.Join(result.Route, " -> "))

	// Atur struktur result menjadi JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to marshal result into JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Output JSON result
	w.Write(resultJSON)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/styles.css", styleHandler)
	http.HandleFunc("/startGame", startGameHandler) // Example handler for starting the game

	// Mulai server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func extractFormData(r *http.Request) (string, string) {
	startName := r.Form.Get("start-article")
	targetName := r.Form.Get("target-article")

	startURLPart := strings.ReplaceAll(startName, " ", "_")
	targetURLPart := strings.ReplaceAll(targetName, " ", "_")

	return startURLPart, targetURLPart
}
