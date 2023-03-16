package api

import (
	"github.com/gin-gonic/gin"
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

func (server *Server) publicRoutes(g *gin.RouterGroup) {
	g.GET("/", server.initPage)
	g.GET("/signup", server.signUpGetHandler)

	g.POST("/login", server.loginUser)
	g.POST("/signup", server.createUser)
}

func (server *Server) privateRoutes(g *gin.RouterGroup) {
	g.GET("/upload", server.uploadGetHandler)
	g.POST("/upload", server.uploadKeywords)
	g.GET("/keywords", server.getKeywords)
	g.GET("/keywords/:id", server.getSearchResultByID)

	g.GET("/logout", server.logoutHandler)

}
