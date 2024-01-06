package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

const (
	chapterHeadingSelector = "#chapter-heading"
	contentSelector        = "body > div.wrap > div > div.site-content > div > div > div > div > div > div > div.c-blog-post > div.entry-content > div > div > div.reading-content > div.text-left > p"
	outputFolder            = "output"
	urlsFile                = "urls.json"
	concurrentWorkers      = 5
)

var wg sync.WaitGroup

func main() {
	c := colly.NewCollector()

	// Read URLs from JSON file
	file, err := os.ReadFile(urlsFile)
	if err != nil {
		log.Fatal("Error reading URLs from file:", err)
	}
	var urls []string
	err = json.Unmarshal(file, &urls)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	// Create a channel to communicate with worker goroutines
	urlsChannel := make(chan string, len(urls))

	// Enqueue URLs to the channel
	for _, url := range urls {
		urlsChannel <- url
	}
	close(urlsChannel)

	// Create a WaitGroup to wait for all worker goroutines to finish
	wg.Add(concurrentWorkers)

	// Start worker goroutines
	for i := 0; i < concurrentWorkers; i++ {
		go func() {
			defer wg.Done()
			// Each worker has its own collector
			wc := c.Clone()

			for url := range urlsChannel {
				var title, content string

				// Set up callback for chapter heading
				wc.OnHTML(chapterHeadingSelector, func(e *colly.HTMLElement) {
					title = e.Text
					log.Println(title)
				})

				// Set up callback for content
				wc.OnHTML(contentSelector, func(e *colly.HTMLElement) {
					trimmedText := strings.TrimSpace(e.Text)
					content += trimmedText + "\n"
				})

				// Visit the URL
				err := wc.Visit(url)
				if err != nil {
					log.Printf("Error visiting URL %s: %v", url, err)
					continue
				}

				// Create the "output" folder if it doesn't exist
				if err := os.MkdirAll(outputFolder, 0755); err != nil {
					log.Fatal(err)
				}

				// Create and open a new file in the "output" folder for writing
				filePath := filepath.Join(outputFolder, title+".txt")
				file, err := os.Create(filePath)
				if err != nil {
					log.Fatal(err)
				}

				// Write the scraped content to the file
				file.WriteString(content)
				file.Close()

				log.Println("Data has been written to", filePath)
			}
		}()
	}

	// Wait for all worker goroutines to finish
	wg.Wait()

	log.Println("All data has been processed.")
}