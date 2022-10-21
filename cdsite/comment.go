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

	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "PUT", "/comment", reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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
	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "POST", "/comment/" + commentID, reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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
	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "DELETE", "/comment/" + commentID, reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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

	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "PUT", "/commentReply", reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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
	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "POST", "/commentReply/" + replyID, reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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
	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "DELETE", "/commentReply/" + replyID, reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
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
	var respBody map[string]interface{}
	isSuccess := CDRequestErrJSON(ctx, sctx, "POST", "/comment/" + commentID + "/respond", reqBody, true, &respBody)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, respBody)
}
