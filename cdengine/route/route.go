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
	r.GET("/course/:cid/comment", func(c *gin.Context) {
		cdengine.GetCourseComment(c, &engineCtx)
	})
	r.GET("/professor/:name", func(c *gin.Context) {
		cdengine.GetProfessor(c, &engineCtx)
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

	// Run CDENGINE
	r.Run(engineCtx.Conf.Host + ":" + strconv.Itoa(engineCtx.Conf.Port))
}
