package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)


// Endpoint "/search"
// Search for a list of course, given search params
func GetSearch(ctx *gin.Context, sctx *Context) {
	// Verify user
	_, isAuth := VerifyUser(ctx, sctx)
	if !isAuth {
		return
	}

	exec := ctx.Query("exec")
	if exec != "true" {
		// Display search page
		ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
			"Title": "Course Search",
			"CourseArray": []cdmodel.Course{}})
		return
	}

	// Retrieve query params
	queryParamsName := []string{
		"query", "is_advanced", "level", "dept", "units"}
	queryParams := make(map[string]interface{})
	for _, name := range queryParamsName {
		queryParams[name] = ctx.Query(name)
	}

	var courseArrResp struct {
		Status		string
		Data		[]cdmodel.Course
	}
	// Execute search
	isSuccess := CDRequest(ctx, sctx, "GET", "/search", queryParams, true, &courseArrResp)
	if !isSuccess {
		return
	}

	// Display course
	ctx.HTML(http.StatusOK, "layout/course_search", gin.H{
		"Title": "Course Search",
		"CurrentQuery": ctx.Query("query"),
		"CourseArray": courseArrResp.Data})
}

