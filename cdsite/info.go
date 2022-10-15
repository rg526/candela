package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// Endpoint "/"
// Homepage
func GetHome(ctx *gin.Context, sctx *Context) {
	// Main content
	ctx.HTML(http.StatusOK, "layout/home", gin.H{
		"Title": "CMU Course List"})
}


// Endpoint "/about"
// About page
func GetAbout(ctx *gin.Context, sctx *Context) {
	// About page
	ctx.HTML(http.StatusOK , "layout/about", gin.H{
		"Title": "About this website"})
}

