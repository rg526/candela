package cdengine

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)

// Endpoint "/course"
// Get course detailed info
func GetCourse(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	var course cdmodel.Course
	cid_query := ctx.Param("cid")
	cid, err := strconv.Atoi(cid_query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Query DB
	stmtCourse, err := ectx.DB.Prepare("SELECT cid, name, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount FROM course WHERE cid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	err = stmtCourse.QueryRow(cid).Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq, &course.FCEHours, &course.FCETeachingRate, &course.FCECourseRate, &course.FCELevel, &course.FCEStudentCount)
	if err != nil {
		// Course not found
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": course})
}


// Endpoint "/professor"
// Get professor detailed info
func GetProfessor(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	var prof cdmodel.Professor
	prof_name := ctx.Param("name")

	// Query DB
	stmtProf, err := ectx.DB.Prepare("SELECT name, RMPRatingClass, RMPRatingOverall FROM professor WHERE name = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	err = stmtProf.QueryRow(prof_name).Scan(&prof.Name, &prof.RMPRatingClass, &prof.RMPRatingOverall)
	if err != nil {
		// Prof not found
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": prof})
}