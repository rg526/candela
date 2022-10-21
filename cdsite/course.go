package cdsite

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)


// Endpoint "/course"
// Get detailed information about a course
func GetCourse(ctx *gin.Context, sctx *Context) {
	// Verify user
	_, isAuth := VerifyUser(ctx, sctx)
	if !isAuth {
		return
	}

	// Find course ID
	cid := ctx.Param("cid")

	// Send CDAPI request
	var courseResp struct {
		Status		string
		Data		cdmodel.Course
	}
	isSuccess := CDRequest(ctx, sctx, "GET", "/course/" + cid, nil, true, &courseResp)
	if !isSuccess {
		return
	}

	if courseResp.Data.CID == "" {
		ReportError(ctx, http.StatusBadRequest,
			"No such course",
			errors.New("Course " + cid + " not found"))
		return
	}

	// Load FCE
	var fceResp struct {
		Status		string
		Data		cdmodel.FCE
	}
	isSuccess = CDRequest(ctx, sctx, "GET", "/course/" + cid + "/fce", nil, true, &fceResp)
	if !isSuccess {
		return
	}

	// Load prof
	var profResp struct {
		Status		string
		Data		[]cdmodel.Prof
	}
	isSuccess = CDRequest(ctx, sctx, "GET", "/course/" + cid + "/prof", nil, true, &profResp)
	if !isSuccess {
		return
	}

	// Load pages
	var pageResp struct {
		Status		string
		Data		[]cdmodel.Page
	}
	isSuccess = CDRequest(ctx, sctx, "GET", "/course/" + cid + "/page", nil, true, &pageResp)
	if !isSuccess {
		return
	}

	// Load tags
	var tagResp struct {
		Status		string
		Data		[]cdmodel.Tag
	}
	isSuccess = CDRequest(ctx, sctx, "GET", "/course/" + cid + "/tag", nil, true, &tagResp)
	if !isSuccess {
		return
	}

	// Load comments
	var commentResp struct {
		Status		string
		Data		[]cdmodel.Comment
	}
	isSuccess = CDRequest(ctx, sctx, "GET", "/course/" + cid + "/comment", nil, true, &commentResp)
	if !isSuccess {
		return
	}

	// Generate HTML
	ctx.HTML(http.StatusOK, "layout/course_page", gin.H{
		"Title": "Course " + courseResp.Data.CID,
		"Course": courseResp.Data,
		"ProfArray": profResp.Data,
		"FCE": fceResp.Data,
		"PageArray": pageResp.Data,
		"TagArray": tagResp.Data,
		"CommentArray": commentResp.Data})
}
