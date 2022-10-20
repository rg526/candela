package main

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"candela/cdengine"
)


func main() {
	engineCtx := cdengine.InitContext("config.json")

	// Setup routes
	r := gin.Default()
	r.GET("/course/:cid", func(c *gin.Context) {
		cdengine.GetCourse(c, &engineCtx)
	})
	r.GET("/course/:cid/fce", func(c *gin.Context) {
		cdengine.GetCourseFCE(c, &engineCtx)
	})
	r.GET("/course/:cid/prof", func(c *gin.Context) {
		cdengine.GetCourseProf(c, &engineCtx)
	})
	r.GET("/course/:cid/comment", func(c *gin.Context) {
		cdengine.GetCourseComment(c, &engineCtx)
	})
	r.GET("/course/:cid/page", func(c *gin.Context) {
		cdengine.GetCoursePage(c, &engineCtx)
	})
	r.GET("/course/:cid/tag", func(c *gin.Context) {
		cdengine.GetCourseTag(c, &engineCtx)
	})
	r.GET("/auth", func(c *gin.Context) {
		cdengine.GetAuth(c, &engineCtx)
	})
	r.GET("/user", func(c *gin.Context) {
		cdengine.GetUser(c, &engineCtx)
	})
	r.GET("/search", func(c *gin.Context) {
		cdengine.GetSearch(c, &engineCtx)
	})

	// comment
	r.PUT("/comment", func(c *gin.Context) {
		cdengine.PutComment(c, &engineCtx)
	})
	r.POST("/comment/:commentID", func(c *gin.Context) {
		cdengine.PostComment(c, &engineCtx)
	})
	r.DELETE("/comment/:commentID", func(c *gin.Context) {
		cdengine.DeleteComment(c, &engineCtx)
	})

	// commentReply
	r.PUT("/commentReply", func(c *gin.Context) {
		cdengine.PutCommentReply(c, &engineCtx)
	})
	r.POST("/commentReply/:replyID", func(c *gin.Context) {
		cdengine.PostCommentReply(c, &engineCtx)
	})
	r.DELETE("/commentReply/:replyID", func(c *gin.Context) {
		cdengine.DeleteCommentReply(c, &engineCtx)
	})

	// commentResponse
	r.POST("/comment/:commentID/respond", func(c *gin.Context) {
		cdengine.PostCommentResponse(c, &engineCtx)
	})

	// Run CDENGINE
	r.Run(engineCtx.Conf.Host + ":" + strconv.Itoa(engineCtx.Conf.Port))
}
