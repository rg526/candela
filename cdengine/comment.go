package cdengine

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Do query
	stmtComment, err := ectx.DB.Prepare("INSERT INTO comment (cid, uid, content, time, anonymous) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	_, err = stmtComment.Exec(reqBody.CID, user.UID, reqBody.Content, time.Now().Unix(), isAnonymous)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Do query
	stmtComment, err := ectx.DB.Prepare("UPDATE comment SET cid = ?, content = ?, time = ?, anonymous = ? WHERE commentID = ? AND uid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	isAnonymous := 0
	if reqBody.Anonymous {
		isAnonymous = 1
	}
	_, err = stmtComment.Exec(reqBody.CID, reqBody.Content, time.Now().Unix(), isAnonymous, commentID, user.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
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
	stmtComment, err := ectx.DB.Prepare("DELETE FROM comment WHERE commentID = ? AND uid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	_, err = stmtComment.Exec(commentID, user.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}

