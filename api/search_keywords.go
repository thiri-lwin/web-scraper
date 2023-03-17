package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gocolly/colly"
	"github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

type SearchInfo struct {
	Keyword            string
	HTMLCode           string
	NumAds             int32
	NumLinks           int32
	TotalSearchResults string
}

type resultDB struct {
	HTMLCode           string
	NumAds             int32
	NumLinks           int32
	TotalSearchResults string
}

type keywordResponse struct {
	Keyword    string
	UploadedAt time.Time
	ResultID   int
}

func (server *Server) uploadKeywords(ctx *gin.Context) {
	// session := sessions.Default(ctx)
	// userEmail := session.Get(util.Userkey)

	userEmail, err := ctx.Cookie(util.Userkey)
	if err != nil {
		log.Println("failed to get cookie:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusOK)
		return
	}

	dbUser, err := server.store.GetUser(ctx, userEmail)
	if err != nil {
		log.Println("Failed to user info from db:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusInternalServerError)
		return
	}

	if dbUser.ID == 0 {
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusNotFound)
		return
	}

	file, _, err := ctx.Request.FormFile("csvfile")
	if err != nil {
		log.Println("failed to read file:", err)
		renderHTML(ctx, gin.H{"title": "Upload", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName)}, "upload.html", http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println(err.Error())
		renderHTML(ctx, gin.H{"title": "Upload", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName)}, "upload.html", http.StatusUnprocessableEntity)
		return
	}
	if len(records) == 0 {
		log.Println("no record in the file")
		renderHTML(ctx, gin.H{"title": "Upload", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName)}, "upload.html", http.StatusOK)
		return
	}

	// assume no header and all records are keywords that need to search
	resultCh := make(chan SearchInfo)
	var wg sync.WaitGroup
	for _, line := range records {
		for _, keyword := range line {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}
			wg.Add(1)
			go searchKeyword(keyword, resultCh, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var allResults []SearchInfo
	var res []keywordResponse

	for result := range resultCh {
		allResults = append(allResults, result)

		// store results in db
		dbRes := resultDB{
			HTMLCode:           result.HTMLCode,
			NumAds:             result.NumAds,
			NumLinks:           result.NumLinks,
			TotalSearchResults: result.TotalSearchResults,
		}
		data, _ := json.Marshal(dbRes)
		insertedRes, err := server.store.InsertSearchResult(ctx, db.SearchResult{
			UserID:  dbUser.ID,
			Keyword: result.Keyword,
			Results: string(data),
		})
		if err != nil {
			log.Println("Failed to save result in db:", err)
		}
		res = append(res, keywordResponse{
			Keyword:    result.Keyword,
			UploadedAt: insertedRes.CreatedAt,
			ResultID:   insertedRes.ID,
		})
	}

	fmt.Println("len of all results:", len(allResults))
	renderHTML(ctx, gin.H{"title": "Upload", "keywords": res, "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "keyword_page": true}, "upload.html", http.StatusOK)
}

func searchKeyword(keyword string, resultCh chan SearchInfo, wg *sync.WaitGroup) {
	// fmt.Println("Keyword:", keyword)
	defer (*wg).Done()

	// Create a new collector instance
	c := colly.NewCollector()

	// Set user agent to avoid being detected as a bot
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	var results SearchInfo
	results.Keyword = keyword
	keyword = strings.ReplaceAll(keyword, " ", "+")
	// On every page request
	c.OnRequest(func(r *colly.Request) {
		// Print the URL being visited
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("#result-stats", func(e *colly.HTMLElement) {
		// Find the total search results for the keyword
		results.TotalSearchResults = e.Text
		fmt.Println("Search Result:", results.TotalSearchResults)
	})

	c.OnHTML(".ads-ad", func(e *colly.HTMLElement) {
		// Find the number of AdWords advertisers on the page
		results.NumAds++
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Find the number of links on the page
		results.NumLinks++
	})

	// On response received
	c.OnResponse(func(r *colly.Response) {
		// Parse the HTML response
		results.HTMLCode = string(r.Body)
	})

	// Visit the first search result page
	c.Visit(fmt.Sprintf("https://www.google.com/search?q=%s&start=%d", keyword, 0))
	resultCh <- results
}

func (server *Server) uploadGetHandler(ctx *gin.Context) {
	// session := sessions.Default(ctx)
	// userEmail := session.Get(util.Userkey)
	userEmail, err := ctx.Cookie(util.Userkey)
	if err != nil {
		log.Println("failed to get cookie:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusOK)
		return
	}

	dbUser, err := server.store.GetUser(ctx, userEmail)
	if err != nil {
		log.Println("Failed to get user info from db:", err)
		renderHTML(ctx, gin.H{"title": "Upload"}, "upload.html", http.StatusInternalServerError)
		return
	}
	if dbUser.ID == 0 {
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusNotFound)
		return
	}
	renderHTML(ctx, gin.H{"title": "Upload", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName)}, "upload.html", http.StatusOK)
}

func (server *Server) getKeywords(ctx *gin.Context) {
	// session := sessions.Default(ctx)
	// userEmail := session.Get(util.Userkey)
	userEmail, err := ctx.Cookie(util.Userkey)
	if err != nil {
		log.Println("failed to get cookie:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusOK)
		return
	}
	dbUser, err := server.store.GetUser(ctx, fmt.Sprint(userEmail))
	if err != nil {
		log.Println("Failed to get user info from db:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusInternalServerError)
		return
	}
	if dbUser.ID == 0 {
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusNotFound)
		return
	}

	dbResults, err := server.store.GetSearchResultsByUserID(ctx, dbUser.ID)
	if err != nil {
		log.Println("Failed to get results by userID:", err)
		renderHTML(ctx, gin.H{"title": "Keyword List", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "keyword_page": true}, "view_keywords.html", http.StatusOK)
		return
	}
	var res []keywordResponse
	for _, dbRes := range dbResults {
		res = append(res, keywordResponse{
			Keyword:    dbRes.Keyword,
			UploadedAt: dbRes.CreatedAt,
			ResultID:   dbRes.ID,
		})
	}
	renderHTML(ctx, gin.H{"title": "Keyword List", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "keywords": res, "keyword_page": true}, "view_keywords.html", http.StatusOK)

}

func (server *Server) getSearchResultByID(ctx *gin.Context) {
	// session := sessions.Default(ctx)
	// userEmail := session.Get(util.Userkey)
	userEmail, err := ctx.Cookie(util.Userkey)
	if err != nil {
		log.Println("failed to get cookie:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusOK)
		return
	}
	dbUser, err := server.store.GetUser(ctx, fmt.Sprint(userEmail))
	if err != nil {
		log.Println("Failed to get user info from db:", err)
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusInternalServerError)
		return
	}
	if dbUser.ID == 0 {
		renderHTML(ctx, gin.H{"title": "Sign In"}, "index.html", http.StatusNotFound)
		return
	}

	id := ctx.Param("id")
	// id := ctx.Query("id")
	keywordID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Invalid id", err)
		return
	}
	searchResult, err := server.store.GetSearchResultByIDAndUserID(ctx, keywordID, dbUser.ID)
	if err != nil {
		log.Println("Failed to get search result by id:", err)
		renderHTML(ctx, gin.H{"title": "Keyword List", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "keyword_page": true}, "view_keywords.html", http.StatusInternalServerError)
		return
	}

	res := SearchInfo{}
	err = json.Unmarshal([]byte(searchResult.Results), &res)
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		return
	}
	res.Keyword = searchResult.Keyword
	// fmt.Println("res <<<<<<<<", res.NumLinks)

	renderHTML(ctx, gin.H{"title": "Result Information", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "payload": res, "result_page": true}, "details.html", http.StatusOK)

}
