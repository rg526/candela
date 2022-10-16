package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)


// Endpoint "/search"
// Search for a list of course, given search params
func GetSearch(ctx *gin.Context, sctx *Context) {
	exec := ctx.Query("exec")
	if exec != "true" {
		// Display search page
		ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
			"Title": "Course Search",
			"CourseArray": []cdmodel.Course{}})
		return
	}

	var courseArrResp struct {
		Status		string
		Data		[]cdmodel.Course
	}
	// Execute search
	isSuccess := CDRequest(ctx, sctx, "/search", nil, true, &courseArrResp)
	if !isSuccess {
		return
	}

	// Display course
	ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
		"Title": "Course Search",
		"CourseArray": courseArrResp.Data})
}

