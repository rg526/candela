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
	rows, err := ectx.DB.
		Query("SELECT cid, name, description, dept, units, prof, prereq, coreq  FROM course LIMIT ?",
			ectx.Conf.MaxSearchResult)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// Append course to array
	for rows.Next() {
		var course cdmodel.Course
		err = rows.Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		courseArr = append(courseArr, course)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": courseArr})
}
