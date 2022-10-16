package cdengine

import (
	"time"
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
	cid:= ctx.Param("cid")

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


// Endpoint "/course/comment"
// Get comments related to a course
func GetCourseComment(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	cid := ctx.Param("cid")

	// Query DB
	stmtComment, err := ectx.DB.Prepare("SELECT commentID, content, time, anonymous, comment.uid, user.name FROM comment INNER JOIN user ON comment.uid = user.uid WHERE cid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	rows, err := stmtComment.Query(cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	defer rows.Close()

	// Append comments to array
	var commentArr []cdmodel.Comment
	for rows.Next() {
		var comment cdmodel.Comment
		var commentTime, commentUID string
		var commentAnonymous int
		err = rows.Scan(&comment.CommentID, &comment.Content, &commentTime,
			&commentAnonymous, &commentUID, &comment.Author)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "ERROR",
				"Error": "Error: " + err.Error()})
			return
		}
		// Protect UID
		comment.Self = commentUID == user.UID
		comment.Anonymous = commentAnonymous == 1
		comment.Time = commentTime // TODO: convert to readable time

		// Convert time from unix ts (string) into readable string
		commentTimeUnix, err := strconv.ParseInt(commentTime, 10, 64)
		if err != nil {
			comment.Time = err.Error()
		} else {
			comment.Time = time.Unix(commentTimeUnix, 0).String()
		}

		commentArr = append(commentArr, comment)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": commentArr})
}
