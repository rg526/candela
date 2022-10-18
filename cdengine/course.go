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
	cid := ctx.Param("cid")

	// Query DB
	err := ectx.DB.
		QueryRow("SELECT cid, name, description, dept, units, prof, prereq, coreq FROM course WHERE cid = ?",
			cid).
		Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq)

	if err != nil && err != sql.ErrNoRows {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
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

	err := ectx.DB.
		QueryRow("SELECT name, RMPRatingClass, RMPRatingOverall FROM professor WHERE name = ?",
			prof_name).
		Scan(&prof.Name, &prof.RMPRatingClass, &prof.RMPRatingOverall)
	if err != nil && err != sql.ErrNoRows {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": prof})
}


// Endpoint "/course/:cid/fce"
// Get FCE data for a course
func GetCourseFCE(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	var fce cdmodel.FCE
	cid := ctx.Param("cid")

	// Query DB
	err := ectx.DB.
		QueryRow("SELECT cid, hours, teachingRate, courseRate, level, studentCount FROM fce WHERE cid = ?",
			cid).
		Scan(&fce.CID, &fce.Hours, &fce.TeachingRate, &fce.CourseRate, &fce.Level, &fce.StudentCount)

	if err != nil && err != sql.ErrNoRows {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": fce})
}


// Endpoint "/course/:cid/comment"
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
	rows, err := ectx.DB.
		Query(`SELECT
		comment.commentID, comment.content, comment.time, comment.anonymous, comment.uid, comment.score, user_comment.name,
		reply.replyID, reply.content, reply.time, reply.anonymous, reply.uid, user_reply.name
		FROM comment
		LEFT JOIN comment_reply AS reply
		ON comment.commentID = reply.commentID
		INNER JOIN user AS user_comment
		ON comment.uid = user_comment.uid
		LEFT JOIN user AS user_reply
		ON reply.uid = user_reply.uid
		WHERE cid = ?
		ORDER BY comment.score DESC, comment.commentID ASC,
			reply.replyID ASC`,
			cid)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
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
			Score			int
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
			&commentQuery.Anonymous, &commentQuery.UID, &commentQuery.Score,
			&commentQuery.Author,
			&replyQuery.ReplyID, &replyQuery.Content, &replyQuery.Time, &replyQuery.Anonymous, &replyQuery.UID, &replyQuery.Author)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Insert comment
		if (len(commentArr) == 0 ||
			commentArr[len(commentArr) - 1].CommentID != commentQuery.CommentID) {
			var comment cdmodel.Comment
			comment.CommentID = commentQuery.CommentID
			comment.Content = commentQuery.Content
			comment.Self = commentQuery.UID == user.UID
			comment.Score = commentQuery.Score

			// Convert time from unix ts (string) into readable string
			commentTimeUnix, err := strconv.ParseInt(commentQuery.Time, 10, 64)
			if err != nil {
				comment.Time = err.Error()
			} else {
				comment.Time = time.Unix(commentTimeUnix, 0).Format("2006-01-02 15:04")
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
				reply.Time = time.Unix(replyTimeUnix, 0).Format("2006-01-02 15:04")
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

// Endpoint "/course/:cid/comment"
// Get pages related to a course
func GetCoursePage(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	cid := ctx.Param("cid")

	// Query DB
	rows, err := ectx.DB.
		Query(`SELECT title, link, content FROM page
			WHERE cid = ?
			ORDER BY priority DESC, pageID ASC`,
			cid)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// Append comments to array
	var pageArr []cdmodel.Page
	for rows.Next() {
		var page cdmodel.Page
		err = rows.Scan(&page.Title, &page.Link, &page.Content)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		pageArr = append(pageArr, page)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": pageArr})
}
