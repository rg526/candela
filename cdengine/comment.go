package cdengine

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Comments

// Endpoint "/comment" PUT
// Create a new comment
// Request body:
// - CID (string)
// - Content (string)
// - Anonymous (bool)
func PutComment(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get request body
	var reqBody struct {
		CID			string
		Content		string
		Anonymous	bool
	}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	// Do query
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	res, err := ectx.DB.
		Exec("INSERT INTO comment (cid, uid, content, time, anonymous) VALUES (?, ?, ?, ?, ?)",
			reqBody.CID, user.UID, reqBody.Content, time.Now().Unix(), isAnonymous)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"CommentID": id})
}

// Endpoint "/comment" POST
// Update comment cid/content
// Comment uid must match current user's uid
// Request body:
// - CID (string)
// - Content (string)
// - Anonymous (bool)
func PostComment(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get comment ID
	commentID := ctx.Param("commentID")

	// Get request body
	var reqBody struct {
		CID			string
		Content		string
		Anonymous	bool
	}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	// Do query
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	_, err = ectx.DB.
		Exec("UPDATE comment SET cid = ?, content = ?, time = ?, anonymous = ? WHERE commentID = ? AND uid = ?",
			reqBody.CID, reqBody.Content, time.Now().Unix(), isAnonymous, commentID, user.UID)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Endpoint "/comment" DELETE
// Delete a comment (commentID)
// Comment uid must match current user's uid
func DeleteComment(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get comment ID
	commentID := ctx.Param("commentID")

	// Do query
	_, err := ectx.DB.
		Exec("DELETE FROM comment WHERE commentID = ? AND uid = ?",
			commentID, user.UID)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Comment replies

// Endpoint "/commentReply" PUT
// Create a new comment
// Request body:
// - CommentID (int)
// - Content (string)
// - Anonymous (bool)
func PutCommentReply(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get request body
	var reqBody struct {
		CommentID	int
		Content		string
		Anonymous	bool
	}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	// Do query
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	res, err := ectx.DB.
		Exec("INSERT INTO comment_reply (commentID, uid, content, time, anonymous) VALUES (?, ?, ?, ?, ?)",
			reqBody.CommentID, user.UID, reqBody.Content, time.Now().Unix(), isAnonymous)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"ReplyID": id})
}

// Endpoint "/commentReply" POST
// Update comment reply
// Reply uid must match current user's uid
// Request body:
// - CommentID (int)
// - Content (string)
// - Anonymous (bool)
func PostCommentReply(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get comment ID
	replyID := ctx.Param("replyID")

	// Get request body
	var reqBody struct {
		CommentID	int
		Content		string
		Anonymous	bool
	}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	// Do query
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	_, err = ectx.DB.
		Exec("UPDATE comment_reply SET commentID = ?, content = ?, time = ?, anonymous = ? WHERE replyID = ? AND uid = ?",
			reqBody.CommentID, reqBody.Content, time.Now().Unix(), isAnonymous, replyID, user.UID)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Endpoint "/commentReply" DELETE
// Delete a comment reply (replyID)
// Reply uid must match current user's uid
func DeleteCommentReply(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get comment ID
	replyID := ctx.Param("replyID")

	// Do query
	_, err := ectx.DB.
		Exec("DELETE FROM comment_reply WHERE replyID = ? AND uid = ?",
			replyID, user.UID)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}



// Comment responses

// Endpoint "/comment/:commentID/respond"
// Respond to a given comment
// Request body:
// - Like (bool)
func PostCommentResponse(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get comment ID
	commentID := ctx.Param("commentID")

	// Get request body
	var reqBody struct {
		Like		bool
	}

	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	if (reqBody.Like) {
		// Insert row into comment_response
		tx, err := ectx.DB.Begin()
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		defer tx.Rollback()
		_, err = tx.
			Exec(`IF NOT EXISTS (
					SELECT uid FROM comment_response
					WHERE uid = ? AND commentID = ?
				) THEN
					INSERT INTO comment_response (commentID, uid, time)
					VALUES (?, ?, ?);
					UPDATE comment SET score = score + 1
					WHERE commentID = ?;
				END IF`,
					user.UID, commentID,
					commentID, user.UID, time.Now().Unix(),
					commentID)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		err = tx.Commit()
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		// Delete from comment_response
		tx, err := ectx.DB.Begin()
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		defer tx.Rollback()
		_, err = tx.
			Exec(`IF EXISTS (
					SELECT uid FROM comment_response
					WHERE uid = ? AND commentID = ?
				) THEN
					DELETE FROM comment_response
					WHERE commentID = ? AND uid = ?;
					UPDATE comment SET score = score - 1
					WHERE commentID = ?;
				END IF`,
					user.UID, commentID,
					commentID, user.UID,
					commentID)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		err = tx.Commit()
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}
