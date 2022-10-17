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
// Modify a comment
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
