package main

import (
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"candela/cdsite"
)


func main() {
	// Read config file
	conf_content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	var conf cdsite.Config
	err = json.Unmarshal(conf_content, &conf)
	if err != nil {
		log.Fatal("Error: read config file: ", err)
	}

	// Setup session
	r := gin.Default()
	store := cookie.NewStore([]byte(conf.CookieSecret))
	r.Use(sessions.Sessions("candela", store))


	// Setup routes
	r.LoadHTMLGlob("../../cdfrontend/template/**/*.tmpl")
	r.Static("/css", "../../cdfrontend/css")
	r.Static("/js", "../../cdfrontend/js")
	r.GET("/", func(c *gin.Context) {
		cdsite.GetHome(c, conf)
	})
	r.GET("/search", func(c *gin.Context) {
		cdsite.GetSearch(c, conf)
	})
	r.GET("/course", func(c *gin.Context) {
		cdsite.GetCourse(c, conf)
	})
	r.GET("/auth", func(c *gin.Context) {
		cdsite.GetAuth(c, conf)
	})
	r.GET("/authCallback", func(c *gin.Context) {
		cdsite.GetAuthCallback(c, conf)
	})
	r.GET("/logout", func(c *gin.Context) {
		cdsite.GetLogout(c, conf)
	})
	r.GET("/about", func(c *gin.Context) {
		cdsite.GetAbout(c, conf)
	})

	// Run CDSITE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
