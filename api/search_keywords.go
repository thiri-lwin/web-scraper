package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/thiri-lwin/web_scraper/util"
)

type keywordResponse struct {
	Keyword    string
	UploadedAt time.Time
	Status     string
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

	input := util.UploadedFile{
		Records: records,
		UserID:  dbUser.ID,
	}
	renderHTML(ctx, gin.H{"title": "Upload", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "scraping_started": true}, "upload.html", http.StatusOK)

	go func(input util.UploadedFile) {
		server.keywordCh <- input
	}(input)

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
			Status:     dbRes.Status,
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

	res := util.SearchInfo{}
	err = json.Unmarshal([]byte(searchResult.Results), &res)
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		return
	}
	res.Keyword = searchResult.Keyword
	// fmt.Println("res <<<<<<<<", res.NumLinks)

	renderHTML(ctx, gin.H{"title": "Result Information", "username": fmt.Sprintf("%s %s", dbUser.FirstName, dbUser.LastName), "payload": res, "result_page": true}, "details.html", http.StatusOK)

}
