package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/web_scraper/util"
)

type keywordResponse struct {
	Keyword    string
	UploadedAt time.Time
	ResultID   int
}

func (server *Server) getKeywords(ctx *gin.Context) {
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
