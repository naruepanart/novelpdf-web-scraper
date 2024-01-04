package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	chapterHeadingSelector = "#chapter-heading"
	contentSelector        = "body > div.wrap > div > div.site-content > div > div > div > div > div > div > div.c-blog-post > div.entry-content > div > div > div.reading-content > div.text-left > p"
	outputFolder            = "output"
	urlsFile                = "urls.json"
)

func main() {
	c := colly.NewCollector()

	// Read URLs from JSON file
	urls, err := readURLsFromFile(urlsFile)
	if err != nil {
		log.Fatal("Error reading URLs from file:", err)
	}

	// Visit each URL
	for _, url := range urls {
		title, content, err := scrapeContent(c, url)
		if err != nil {
			log.Printf("Error scraping content from URL %s: %v", url, err)
			continue
		}

		// Create the "output" folder if it doesn't exist
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			err := os.Mkdir(outputFolder, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Create and open a new file in the "output" folder for writing
		filePath := filepath.Join(outputFolder, title+".txt")
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write the scraped content to the file
		file.WriteString(content)

		log.Println("Data has been written to", filePath)
	}
}

func scrapeContent(c *colly.Collector, url string) (string, string, error) {
	var title, content string

	// Set up callback for chapter heading
	c.OnHTML(chapterHeadingSelector, func(e *colly.HTMLElement) {
		title = e.Text
		log.Println(title)
	})

	// Set up callback for content
	c.OnHTML(contentSelector, func(e *colly.HTMLElement) {
		trimmedText := strings.TrimSpace(e.Text)
		content += trimmedText + "\n"
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		return "", "", err
	}

	return title, content, nil
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var urls []string
	err = json.Unmarshal(file, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}