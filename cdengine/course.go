package cdengine

import (
	"time"
	"strconv"
	"net/http"
	"database/sql"
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
	stmtComment, err := ectx.DB.Prepare(`SELECT
		comment.commentID, comment.content, comment.time, comment.anonymous, comment.uid, user_comment.name,
		reply.replyID, reply.content, reply.time, reply.anonymous, reply.uid, user_reply.name
		FROM comment
		LEFT JOIN comment_reply AS reply
		ON comment.commentID = reply.commentID
		INNER JOIN user AS user_comment
		ON comment.uid = user_comment.uid
		LEFT JOIN user AS user_reply
		ON reply.uid = user_reply.uid
		WHERE cid = ?`)
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
		var commentQuery struct {
			CommentID		int
			Content			string
			Time			string
			Anonymous		int
			UID				string
			Author			string
		}
		var replyQuery struct {
			ReplyID			sql.NullInt32
			Content			sql.NullString
			Time			sql.NullString
			Anonymous		sql.NullInt32
			UID				sql.NullString
			Author			sql.NullString
		}

		err = rows.Scan(&commentQuery.CommentID, &commentQuery.Content, &commentQuery.Time,
			&commentQuery.Anonymous, &commentQuery.UID, &commentQuery.Author,
			&replyQuery.ReplyID, &replyQuery.Content, &replyQuery.Time, &replyQuery.Anonymous, &replyQuery.UID, &replyQuery.Author)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "ERROR",
				"Error": "Error: " + err.Error()})
			return
		}

		// Insert comment
		if (len(commentArr) == 0 ||
			commentArr[len(commentArr) - 1].CommentID != commentQuery.CommentID) {
			var comment cdmodel.Comment
			comment.CommentID = commentQuery.CommentID
			comment.Content = commentQuery.Content
			comment.Self = commentQuery.UID == user.UID

			// Convert time from unix ts (string) into readable string
			commentTimeUnix, err := strconv.ParseInt(commentQuery.Time, 10, 64)
			if err != nil {
				comment.Time = err.Error()
			} else {
				comment.Time = time.Unix(commentTimeUnix, 0).String()
			}

			// Protect Anonymous
			if (commentQuery.Anonymous == 1) {
				comment.Author = "Anonymous"
			} else {
				comment.Author = commentQuery.Author
			}

			// Insert to commentArr
			commentArr = append(commentArr, comment)
		}

		// Insert reply
		if (replyQuery.ReplyID.Valid) {
			var reply cdmodel.CommentReply
			reply.ReplyID = int(replyQuery.ReplyID.Int32)
			reply.Content = replyQuery.Content.String
			reply.Self = replyQuery.UID.String == user.UID

			// Convert time
			replyTimeUnix, err := strconv.ParseInt(replyQuery.Time.String, 10, 64)
			if err != nil {
				reply.Time = err.Error()
			} else {
				reply.Time = time.Unix(replyTimeUnix, 0).String()
			}

			// Protect Anonymous
			if (replyQuery.Anonymous.Int32 == 1) {
				reply.Author = "Anonymous"
			} else {
				reply.Author = replyQuery.Author.String
			}

			// Insert into last comment
			replies := &commentArr[len(commentArr) - 1].Replies
			*replies = append(*replies, reply)
		}

	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": commentArr})
}
