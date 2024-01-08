package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

const (
	chapterHeadingSelector = "#chapter-heading"
	contentSelector        = "body > div.wrap > div > div.site-content > div > div > div > div > div > div > div.c-blog-post > div.entry-content > div > div > div.reading-content > div.text-left > p"
	outputFolder           = "output"
	urlsFile               = "urls.json"
	concurrentWorkers      = 10
	memoryLimit            = 100 * 1024 * 1024 // 100 MB
)

var wg sync.WaitGroup

func main() {
	debug.SetMemoryLimit(memoryLimit)

	c := colly.NewCollector()

	urls := readURLsFromJSON(urlsFile)

	urlsChannel := make(chan string, len(urls))
	for _, url := range urls {
		urlsChannel <- url
	}
	close(urlsChannel)

	wg.Add(concurrentWorkers)

	for i := 0; i < concurrentWorkers; i++ {
		go worker(c.Clone(), urlsChannel)
	}

	wg.Wait()

	log.Println("All data has been processed.")
}

func readURLsFromJSON(filename string) []string {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading URLs from file:", err)
	}

	var urls []string
	err = json.Unmarshal(file, &urls)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	return urls
}

func worker(wc *colly.Collector, urlsChannel <-chan string) {
	defer wg.Done()

	var title, content string

	wc.OnHTML(chapterHeadingSelector, func(e *colly.HTMLElement) {
		title = e.Text
		log.Println(title)
	})

	wc.OnHTML(contentSelector, func(e *colly.HTMLElement) {
		trimmedText := strings.TrimSpace(e.Text)
		content += trimmedText + "\n"
	})

	for url := range urlsChannel {
		title, content = "", ""

		if err := wc.Visit(url); err != nil {
			log.Printf("Error visiting URL %s: %v", url, err)
			continue
		}

		if err := os.MkdirAll(outputFolder, 0755); err != nil {
			log.Fatal(err)
		}

		sanitizedTitle := sanitizeTitle(title)
		filePath := filepath.Join(outputFolder, sanitizedTitle+".txt")

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString(content)
		file.Close()
		
		log.Println("Data has been written to", filePath)
	}
}

func sanitizeTitle(title string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		title = strings.ReplaceAll(title, char, "_")
	}

	return strings.TrimSpace(title)
}
