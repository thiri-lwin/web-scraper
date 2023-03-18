package util

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type UploadedFile struct {
	Records [][]string
	UserID  int
}

type SearchInfo struct {
	Keyword            string
	HTMLCode           string
	NumAds             int32
	NumLinks           int32
	TotalSearchResults string
	Status             string // "pending", "scraping", "complete", or "error"
}

type ResultDB struct {
	HTMLCode           string
	NumAds             int32
	NumLinks           int32
	TotalSearchResults string
}

func SearchKeyword(keywordInfo SearchInfo, resultCh chan SearchInfo, wg *sync.WaitGroup) {
	// fmt.Println("Keyword:", keyword)
	defer (*wg).Done()

	keywordInfo.Status = "scraping"

	// Create a new collector instance
	c := colly.NewCollector()

	// Set user agent to avoid being detected as a bot
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	keyword := strings.ReplaceAll(keywordInfo.Keyword, " ", "+")
	// On every page request
	c.OnRequest(func(r *colly.Request) {
		// Print the URL being visited
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("#result-stats", func(e *colly.HTMLElement) {
		// Find the total search results for the keyword
		keywordInfo.TotalSearchResults = e.Text
		fmt.Println("Search Result:", keywordInfo.TotalSearchResults)
	})

	c.OnHTML(".ads-ad", func(e *colly.HTMLElement) {
		// Find the number of AdWords advertisers on the page
		keywordInfo.NumAds++
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Find the number of links on the page
		keywordInfo.NumLinks++
	})

	// On response received
	c.OnResponse(func(r *colly.Response) {
		// Parse the HTML response
		keywordInfo.HTMLCode = string(r.Body)
	})

	// Visit the first search result page
	c.Visit(fmt.Sprintf("https://www.google.com/search?q=%s&start=%d", keyword, 0))
	if keywordInfo.HTMLCode != "" {
		keywordInfo.Status = "completed"
	} else {
		keywordInfo.Status = "error"
	}
	resultCh <- keywordInfo
}
