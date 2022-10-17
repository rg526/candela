package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// Endpoint "/comment" (PUT)
// Create a new comment
func PutComment(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	isSuccess := CDRequest(ctx, sctx, "PUT", "/comment", reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}

// Endpoint "/comment" (POST)
// Modify a comment
func PostComment(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	commentID := ctx.Param("commentID")
	isSuccess := CDRequest(ctx, sctx, "POST", "/comment/" + commentID, reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Endpoint "/comment" (DELETE)
// Delete a comment
func DeleteComment(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	commentID := ctx.Param("commentID")
	isSuccess := CDRequest(ctx, sctx, "DELETE", "/comment/" + commentID, reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Endpoint "/commentReply" (PUT)
// Create a new comment reply
func PutCommentReply(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	isSuccess := CDRequest(ctx, sctx, "PUT", "/commentReply", reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}

// Endpoint "/commentReply" (POST)
// Modify a comment reply
func PostCommentReply(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	replyID := ctx.Param("replyID")
	isSuccess := CDRequest(ctx, sctx, "POST", "/commentReply/" + replyID, reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}


// Endpoint "/commentReply" (DELETE)
// Delete a comment reply
func DeleteCommentReply(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	replyID := ctx.Param("replyID")
	isSuccess := CDRequest(ctx, sctx, "DELETE", "/commentReply/" + replyID, reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}

// Endpoint "/comment/:commentID/respond" (POST)
// Respond to a comment
func PostCommentResponse(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	commentID := ctx.Param("commentID")
	isSuccess := CDRequest(ctx, sctx, "POST", "/comment/" + commentID + "/respond", reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}
