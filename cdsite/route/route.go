package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"candela/cdsite"
)


func main() {
	siteCtx := cdsite.InitContext("config.json")

	// Setup session
	r := gin.Default()
	store := cookie.NewStore([]byte(siteCtx.Conf.CookieSecret))
	r.Use(sessions.Sessions("candela", store))

	// Setup routes
	r.LoadHTMLGlob("../../cdfrontend/template/**/*.tmpl")
	r.Static("/css", "../../cdfrontend/css")
	r.Static("/js", "../../cdfrontend/js")
	r.StaticFile("/robots.txt", "../../cdfrontend/resource/robots.txt")
	r.StaticFile("/manifest.json", "../../cdfrontend/resource/manifest.json")
	r.GET("/", func(c *gin.Context) {
		cdsite.GetHome(c, &siteCtx)
	})
	r.GET("/search", func(c *gin.Context) {
		cdsite.GetSearch(c, &siteCtx)
	})
	r.GET("/course/:cid", func(c *gin.Context) {
		cdsite.GetCourse(c, &siteCtx)
	})
	r.GET("/auth", func(c *gin.Context) {
		cdsite.GetAuth(c, &siteCtx)
	})
	r.GET("/auth/callback", func(c *gin.Context) {
		cdsite.GetAuthCallback(c, &siteCtx)
	})
	r.GET("/logout", func(c *gin.Context) {
		cdsite.GetLogout(c, &siteCtx)
	})
	r.GET("/about", func(c *gin.Context) {
		cdsite.GetAbout(c, &siteCtx)
	})

	// comment
	r.PUT("/comment", func(c *gin.Context) {
		cdsite.PutComment(c, &siteCtx)
	})
	r.POST("/comment/:commentID", func(c *gin.Context) {
		cdsite.PostComment(c, &siteCtx)
	})
	r.DELETE("/comment/:commentID", func(c *gin.Context) {
		cdsite.DeleteComment(c, &siteCtx)
	})
	// commentReply
	r.PUT("/commentReply", func(c *gin.Context) {
		cdsite.PutCommentReply(c, &siteCtx)
	})
	r.POST("/commentReply/:replyID", func(c *gin.Context) {
		cdsite.PostCommentReply(c, &siteCtx)
	})
	r.DELETE("/commentReply/:replyID", func(c *gin.Context) {
		cdsite.DeleteCommentReply(c, &siteCtx)
	})
	// commentResponse
	r.POST("/comment/:commentID/respond", func(c *gin.Context) {
		cdsite.PostCommentResponse(c, &siteCtx)
	})
	// tag
	r.PUT("/tag", func(c *gin.Context) {
		cdsite.PutTag(c, &siteCtx)
	})

	// Run CDSITE
	r.Run(siteCtx.Conf.Host + ":" + strconv.Itoa(siteCtx.Conf.Port))
}
