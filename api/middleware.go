package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/web_scraper/util"

	"net/http"
)

func authRequired(c *gin.Context) {
	// session := sessions.Default(c)
	// user := session.Get("user")
	userEmail, err := c.Cookie(util.Userkey)
	if userEmail == "" || err != nil {
		log.Println("User not logged in")
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}
	c.Next()
}
