package cdsite

import (
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)

// Endpoint "/search"
// Search for a list of course, given search params
func GetSearch(ctx *gin.Context, sctx *Context) {
	// Main content
	ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
		"Title": "Course Search"})
}


// Endpoint "/course"
// Get detailed information about a course
func GetCourse(ctx *gin.Context, sctx *Context) {
	// Find course ID
	cid := ctx.Param("cid")
	courseUrl := "/course/" + cid

	// Send CDAPI request
	var courseResp struct {
		Status		string
		Data		cdmodel.Course
	}
	CDRequest(ctx, sctx, courseUrl, url.Values{}, true, &courseResp)

	// Fetch prof struct
	profName := strings.Split(courseResp.Data.Prof, ";")
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
		profUrl := sctx.Conf.CDAPIUrl + "/professor/" + name

		// Do request
		req, err := http.NewRequest("GET", profUrl, nil)
		if err != nil {
			res, err := sctx.Client.Do(req)
			if err == nil && res.StatusCode == http.StatusOK {
				_ = json.NewDecoder(res.Body).Decode(&prof)
			}
		}
		profArr = append(profArr, prof.Data)
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "layout/course_page", gin.H{
		"Title": "Course " + strconv.Itoa(courseResp.Data.CID),
		"Course": courseResp.Data,
		"ProfArray": profArr})
}
