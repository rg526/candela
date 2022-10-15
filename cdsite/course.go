package cdsite

import (

	"strconv"
	"strings"
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	"candela/cdmodel"
)

func GetSearch(ctx *gin.Context, sctx *Context) {
	// AUTH REQUIRED
	session := sessions.Default(ctx)
	token := session.Get("token")
	if token == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Main content
	ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
		"Title": "Course Search"})
}

func GetCourse(ctx *gin.Context, sctx *Context) {
	// AUTH REQUIRED
	token, _, isAuth := VerifyUserFromSession(ctx, sctx)
	if !isAuth {
		return
	}

	// Find course ID
	cid := ctx.Param("cid")
	courseVal := url.Values{}
	courseVal.Add("token", token)
	courseUrl := sctx.Conf.CDAPIUrl + "course/" + cid + "/?" + courseVal.Encode()


	// Send CDAPI request
	var course struct {
		Status		string
		Data		cdmodel.Course
	}
	res, err := http.Get(courseUrl)
	if err != nil {
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Unsuccessful connection to CDEngine."})
		return
	}
	if res.StatusCode != http.StatusOK {
		var msg map[string]interface{}
		json.NewDecoder(res.Body).Decode(&msg)
		ctx.HTML(http.StatusBadRequest, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Invalid Request",
			"ErrorDescription": msg["Error"].(string)})
		return
	}

	err = json.NewDecoder(res.Body).Decode(&course)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Decode error: " + err.Error()})
		return
	}

	// Fetch prof struct
	profName := strings.Split(course.Data.Prof, ";")
	var profArr []cdmodel.Professor
	for _, name := range(profName) {
		// Default data
		var prof struct {
			Status string
			Data cdmodel.Professor
		}
		prof.Data.Name = name
		prof.Data.RMPRatingClass = "Unknown"
		prof.Data.RMPRatingOverall = -1.0

		// Build URL
		profVal := url.Values{}
		profVal.Add("token", token)
		profUrl := sctx.Conf.CDAPIUrl + "professor/" + name + "/?" + profVal.Encode()

		// Do request
		res, err = http.Get(profUrl)
		if err == nil && res.StatusCode == http.StatusOK {
			_ = json.NewDecoder(res.Body).Decode(&prof)
		}
		profArr = append(profArr, prof.Data)
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "layout/course_page", gin.H{
		"Title": "Course " + strconv.Itoa(course.Data.CID),
		"Course": course.Data,
		"ProfArray": profArr})
}
