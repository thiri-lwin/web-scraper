package util

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestSearchKeyword(t *testing.T) {
	resultCh := make(chan SearchInfo)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	keyword := "test keyword"
	arg := SearchInfo{
		Keyword: keyword,
		Status:  "pending",
	}
	go SearchKeyword(arg, resultCh, wg)
	results := <-resultCh
	wg.Wait()
	// expectedTotalSearchResults := "About 1,000,000 results"
	// if results.TotalSearchResults != expectedTotalSearchResults {
	// 	t.Errorf("Unexpected total search results. Expected: %s, but got: %s", expectedTotalSearchResults, results.TotalSearchResults)
	// }
	if results.TotalSearchResults == "" {
		t.Error("Unexpected total search results. Total search results should not be empty.")
	}
	expectedTitle := fmt.Sprintf("<title>%s - ", keyword)
	if !(strings.Index(results.HTMLCode, expectedTitle) > 0) {
		t.Error("Unexpected HTML Content.")
	}
	status := "completed"
	if results.Status != status {
		t.Errorf("Unexpected status. Expected: %s, but got: %s", status, results.Status)
	}
}
