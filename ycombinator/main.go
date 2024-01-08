package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	var hrefs []string

	c.OnHTML("td:nth-child(3) > span", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		hrefs = append(hrefs, link)
		fmt.Printf("Link: %s\n", link)
	})

	c.Visit("https://news.ycombinator.com/")

	file, err := os.Create("urls.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(hrefs); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON data saved to urls.json")
}