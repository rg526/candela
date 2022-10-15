package cdsite

import (
	"strconv"
	"strings"
	"net/http"
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

	// Send CDAPI request
	var courseResp struct {
		Status		string
		Data		cdmodel.Course
	}
	isSuccess := CDRequest(ctx, sctx, "/course/" + cid, nil, true, &courseResp)
	if !isSuccess {
		return
	}

	if courseResp.Data.CID == 0 {
		ctx.HTML(http.StatusBadRequest, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "No such course",
			"ErrorDescription": "Course " + cid + " not found"})
		return
	}

	// Fetch prof struct
	profName := strings.Split(courseResp.Data.Prof, ";")
	var profArr []cdmodel.Professor
	for _, name := range(profName) {
		var profResp struct {
			Status string
			Data cdmodel.Professor
		}
		isSuccess := CDRequest(ctx, sctx, "/professor/" + name, nil, true, &profResp)
		if !isSuccess {
			return
		}

		// Default data
		if profResp.Data.Name == "" {
			profResp.Data.Name = name
			profResp.Data.RMPRatingClass = "Unknown"
			profResp.Data.RMPRatingOverall = -1.0
		}
		profArr = append(profArr, profResp.Data)
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "layout/course_page", gin.H{
		"Title": "Course " + strconv.Itoa(courseResp.Data.CID),
		"Course": courseResp.Data,
		"ProfArray": profArr})
}
