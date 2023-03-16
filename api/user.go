package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	db "github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

type createUserRequest struct {
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
	Email     string `form:"email" binding:"required"`
	Password  string `form:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error in binding request:", err)
		renderHTML(ctx, gin.H{"title": "Sign Up", "content": "Something went wrong."}, "signup.html", http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Println("Error in password hashing:", err)
		renderHTML(ctx, gin.H{"title": "Sign Up", "content": "Something went wrong."}, "signup.html", http.StatusInternalServerError)
		return
	}

	arg := db.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				renderHTML(ctx, gin.H{"title": "Sign Up", "content": "Email is already registered."}, "signup.html", http.StatusBadRequest)
				return
			}
		}
		renderHTML(ctx, gin.H{"title": "Sign Up", "content": "Something went wrong."}, "signup.html", http.StatusInternalServerError)
		return
	}

	// rsp := userResponse{
	// 	FirstName: user.FirstName,
	// 	LastName:  user.LastName,
	// 	Email:     user.Email,
	// 	CreatedAt: user.CreatedAt,
	// }
	ctx.Redirect(http.StatusFound, "/")
}

type loginUserRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	user, _ := ctx.Cookie(util.Userkey)
	if user != "" {
		ctx.SetCookie(util.Userkey, "", -1, "", "", false, true)

	}

	var req loginUserRequest
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding input json :", err.Error())
		renderHTML(ctx, gin.H{"title": "Sign In", "content": "Something went wrong. Please try again later."}, "index.html", http.StatusBadRequest)
		return
	}

	dbUser, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			renderHTML(ctx, gin.H{"title": "Sign In", "content": "Incorrect email or password."}, "index.html", http.StatusNotFound)
			return
		}
		renderHTML(ctx, gin.H{"title": "Sign In", "content": "Something went wrong. Please try again later."}, "index.html", http.StatusInternalServerError)
		return
	}

	if dbUser.ID == 0 {
		renderHTML(ctx, gin.H{"title": "Sign In", "content": "Incorrect email or password."}, "index.html", http.StatusNotFound)
		return
	}

	err = util.CheckPassword(req.Password, dbUser.HashedPassword)
	if err != nil {
		log.Println("Error in checking password:", err)
		renderHTML(ctx, gin.H{"title": "Sign In", "content": "Incorrect email or password."}, "index.html", http.StatusUnauthorized)
		return
	}

	ctx.SetCookie(util.Userkey, req.Email, 3600, "", "", false, true)

	// fmt.Println("session >>>>>>>>>>>", session.Get(util.Userkey))
	ctx.Redirect(http.StatusMovedPermanently, "/keywords")

}
