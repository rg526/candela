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
	r.GET("/course", func(c *gin.Context) {
		cdengine.GetCourse(c, &engineCtx)
	})
	r.GET("/professor", func(c *gin.Context) {
		cdengine.GetProfessor(c, &engineCtx)
	})
	r.GET("/auth", func(c *gin.Context) {
		cdengine.GetAuth(c, &engineCtx)
	})
	r.GET("/user", func(c *gin.Context) {
		cdengine.GetUser(c, &engineCtx)
	})

	// Run CDENGINE
	r.Run(engineCtx.Conf.Host + ":" + strconv.Itoa(engineCtx.Conf.Port))
}
