package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// Endpoint "/"
// Homepage
func GetHome(ctx *gin.Context, sctx *Context) {
	// Verify user
	_, isAuth := VerifyUser(ctx, sctx)
	if !isAuth {
		return
	}

	// Main content
	ctx.HTML(http.StatusOK, "layout/home", gin.H{
		"Title": "CMU Course List"})
}


// Endpoint "/about"
// About page
func GetAbout(ctx *gin.Context, sctx *Context) {
	// Verify user
	_, isAuth := VerifyUser(ctx, sctx)
	if !isAuth {
		return
	}

	// About page
	ctx.HTML(http.StatusOK , "layout/about", gin.H{
		"Title": "About this website"})
}

