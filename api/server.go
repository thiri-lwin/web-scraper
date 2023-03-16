package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	db "github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

type Server struct {
	store        *db.Store
	router       *gin.Engine
	templatePath string
	assetsPath   string
	cssPath      string
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store, templatePath, assetsPath, cssPath string) *Server {
	server := &Server{store: store, templatePath: templatePath, assetsPath: assetsPath, cssPath: cssPath}

	// set routes
	server.InitRoutes()
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) InitRoutes() {
	router := gin.Default()
	router.LoadHTMLGlob(server.templatePath)
	router.Static(server.assetsPath, "."+server.assetsPath)
	router.Static(server.cssPath, "./"+server.cssPath)
	router.Use(sessions.Sessions("session", cookie.NewStore(util.Secret)))

	router.GET("/", server.initPage)
	router.GET("/signup", server.signUpGetHandler)
	server.router = router
}

func (server *Server) initPage(c *gin.Context) {
	renderHTML(c, gin.H{"title": "Sign In"}, "index.html", http.StatusOK)
}

func (server *Server) signUpGetHandler(c *gin.Context) {
	renderHTML(c, gin.H{"title": "Sign Up"}, "signup.html", http.StatusOK)
}

func renderHTML(c *gin.Context, data gin.H, templateName string, code int) {
	c.HTML(code, templateName, data)
}
