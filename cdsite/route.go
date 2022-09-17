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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

type CourseResponse struct {
	Status		string
	Data		cdmodel.Course
}

type ProfessorResponse struct {
	Status		string
	Data        cdmodel.Professor
}

func getHome(ctx *gin.Context, conf config) {
	// AUTH REQUIRED
	session := sessions.Default(ctx)
	if session.Get("token") == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Main content
	ctx.HTML(http.StatusOK, "home_page.tmpl", gin.H{})
}

func getCourse(ctx *gin.Context, conf config) {
	// AUTH REQUIRED
	session := sessions.Default(ctx)
	if session.Get("token") == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

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

func getAuth(ctx *gin.Context, conf config) {
	// Redirect to Google
	authVal := url.Values{}
	authVal.Add("client_id", conf.OAuth2ClientID)
	authVal.Add("redirect_uri", conf.OAuth2RedirectURI)
	authVal.Add("response_type", "code")
	authVal.Add("scope", conf.OAuth2Scope)
	authVal.Add("hd", "andrew.cmu.edu")
	authVal.Add("prompt", "select_account")
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth?" + authVal.Encode()

	ctx.Redirect(http.StatusMovedPermanently, authUrl)
}

type AuthResponse struct {
	AccessToken		string		`json:"access_token"`
}

type UserResponse struct {
	Status			string
	Token			string
}

func getAuthCallback(ctx *gin.Context, conf config) {
	// Get authCode
	if ctx.Query("error") != "" {
		ctx.HTML(http.StatusOK, "error_page.tmpl", gin.H{
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Google returns " + ctx.Query("error") + "."})
		return
	}
	authCode := ctx.Query("code")

	// Verify using CDEngine
	authVal := url.Values{}
	authVal.Add("code", authCode)
	res, err := http.Get(conf.CDAPIUrl + "auth?" + authVal.Encode())
	if err != nil {
		ctx.HTML(http.StatusOK, "error_page.tmpl", gin.H{
				"ErrorTitle": "Service Error",
				"ErrorDescription": "Unsuccessful connection to CDEngine."})
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.HTML(http.StatusOK, "error_page.tmpl", gin.H{
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Invalid request."})
		return
	}
	var user UserResponse
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		ctx.HTML(http.StatusOK, "error_page.tmpl", gin.H{
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Error decoding json" + err.Error() + "."})
		return
	}

	// Set cookie
	session := sessions.Default(ctx)
	session.Set("token", user.Token)
	session.Save()

	// Page
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}


type config struct {
	Host				string
	Port				int
	CDAPIUrl			string
	CookieSecret		string
	OAuth2ClientID		string
	OAuth2ClientSecret	string
	OAuth2Scope			string
	OAuth2RedirectURI	string
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

	// Setup session
	r := gin.Default()
	store := cookie.NewStore([]byte(conf.CookieSecret))
	r.Use(sessions.Sessions("candela", store))


	// Setup routes
	r.LoadHTMLGlob("../cdfrontend/*.tmpl")
	r.Static("/css", "../cdfrontend/css")
	r.Static("/js", "../cdfrontend/js")
	r.StaticFile("/about", "../cdfrontend/about.tmpl")
	r.GET("/", func(c *gin.Context) {
		getHome(c, conf)
	})
	r.GET("/course", func(c *gin.Context) {
		getCourse(c, conf)
	})
	r.GET("/auth", func(c *gin.Context) {
		getAuth(c, conf)
	})
	r.GET("/authCallback", func(c *gin.Context) {
		getAuthCallback(c, conf)
	})

	// Run CDSITE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
