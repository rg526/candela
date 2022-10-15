package cdengine

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)

// Endpoint "/search"
// Search for courses
func GetSearch(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	var courseArr []cdmodel.Course

	// Query DB
	stmtCourse, err := ectx.DB.Prepare("SELECT cid, name, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount FROM course LIMIT ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	rows, err := stmtCourse.Query(ectx.Conf.MaxSearchResult)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	defer rows.Close()

	// Append course to array
	for rows.Next() {
		var course cdmodel.Course
		err = rows.Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq, &course.FCEHours, &course.FCETeachingRate, &course.FCECourseRate, &course.FCELevel, &course.FCEStudentCount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "ERROR",
				"Error": "Error: " + err.Error()})
			return
		}
		courseArr = append(courseArr, course)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": courseArr})
}
