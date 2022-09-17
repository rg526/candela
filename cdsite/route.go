package main

import (
	"log"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/url"

	"candela/cdmodel"

	"github.com/gin-gonic/gin"
	"net/http"
)

type CourseResponse struct {
	Status		string
	Data		cdmodel.Course
}

type ProfessorResponse struct {
	Status		string
	Data		cdmodel.Professor
}

func getCourse(ctx *gin.Context, conf config) {
	// Find course ID
	var course CourseResponse
	courseVal := url.Values{}
	courseVal.Add("cid", ctx.Query("cid"))
	courseUrl := conf.CDAPIUrl + "course?" + courseVal.Encode()


	// Send CDAPI request
	res, err := http.Get(courseUrl)
	if err != nil {
		ctx.HTML(http.StatusServiceUnavailable, "error_page.tmpl", gin.H{
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Unsuccessful connection to CDEngine."})
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.HTML(http.StatusServiceUnavailable, "error_page.tmpl", gin.H{
			"ErrorTitle": "Invalid Request",
			"ErrorDescription": "Your request is not valid."})
		return
	}
	err = json.NewDecoder(res.Body).Decode(&course)
	if err != nil {
		ctx.HTML(http.StatusServiceUnavailable, "error_page.tmpl", gin.H{
			"ErrorTitle": "Service Error",
			"ErrorDescription": "CDEngine returns invalid response."})
		return
	}

	// Fetch prof struct
	profName := strings.Split(course.Data.Prof, ";")
	var profArr []cdmodel.Professor
	for _, name := range(profName) {
		// Default data
		var prof ProfessorResponse
		prof.Data.Name = name
		prof.Data.RMPRatingClass = "Unknown"
		prof.Data.RMPRatingOverall = -1.0

		// Build URL
		profVal := url.Values{}
		profVal.Add("name", name)
		profUrl := conf.CDAPIUrl + "professor?" + profVal.Encode()

		// Do request
		res, err = http.Get(profUrl)
		if err == nil && res.StatusCode == http.StatusOK {
			_ = json.NewDecoder(res.Body).Decode(&prof)
		}
		profArr = append(profArr, prof.Data)
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "course_page.tmpl", gin.H{
		"Course": course.Data,
		"ProfArray": profArr})
}


type config struct {
	Host			string
	Port			int
	CDAPIUrl		string
}

func main() {
	// Read config file
	conf_content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	var conf config
	err = json.Unmarshal(conf_content, &conf)
	if err != nil {
		log.Fatal("Error: read config file: ", err)
	}

	// Setup routes
	r := gin.Default()
	r.LoadHTMLGlob("../cdfrontend/*.tmpl")
	r.Static("/css", "../cdfrontend/css")
	r.Static("/js", "../cdfrontend/js")
	r.StaticFile("/about", "../cdfrontend/about.tmpl")
	r.GET("/course", func(c *gin.Context) {
		getCourse(c, conf)
	})

	// Run CDSITE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
