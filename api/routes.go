package api

import "github.com/gin-gonic/gin"

func (server *Server) publicRoutes(g *gin.RouterGroup) {
	g.GET("/", server.initPage)
	g.GET("/signup", server.signUpGetHandler)

	g.POST("/login", server.loginUser)
	g.POST("/signup", server.createUser)
}

func (server *Server) privateRoutes(g *gin.RouterGroup) {

}
