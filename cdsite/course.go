package cdsite

import (
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)

// Endpoint "/search"
// Search for a list of course, given search params
func GetSearch(ctx *gin.Context, sctx *Context) {
	// AUTH REQUIRED
	_, _, isAuth := VerifyUserFromSession(ctx, sctx)
	if !isAuth {
		return
	}

	// Main content
	ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
		"Title": "Course Search"})
}


// Endpoint "/course"
// Get detailed information about a course
func GetCourse(ctx *gin.Context, sctx *Context) {
	// AUTH REQUIRED
	token, _, isAuth := VerifyUserFromSession(ctx, sctx)
	if !isAuth {
		return
	}

	// Find course ID
	cid := ctx.Param("cid")
	courseUrl := sctx.Conf.CDAPIUrl + "course/" + cid


	// Send CDAPI request
	var course struct {
		Status		string
		Data		cdmodel.Course
	}
	req, err := http.NewRequest("GET", courseUrl, nil)
	if err != nil {
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Connection error: " + err.Error()})
		return
	}
	req.Header.Add("Authorization", token)
	res, err := sctx.Client.Do(req)
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
		profUrl := sctx.Conf.CDAPIUrl + "professor/" + name

		// Do request
		req, err := http.NewRequest("GET", profUrl, nil)
		if err != nil {
			res, err = sctx.Client.Do(req)
			if err == nil && res.StatusCode == http.StatusOK {
				_ = json.NewDecoder(res.Body).Decode(&prof)
			}
		}
		profArr = append(profArr, prof.Data)
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "layout/course_page", gin.H{
		"Title": "Course " + strconv.Itoa(course.Data.CID),
		"Course": course.Data,
		"ProfArray": profArr})
}
