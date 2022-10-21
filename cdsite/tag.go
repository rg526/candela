package cdsite

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// Endpoint "/tag" (PUT)
// Create a new tag 
func PutTag(ctx *gin.Context, sctx *Context) {
	var reqBody map[string]interface{}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	isSuccess := CDRequestErrJSON(ctx, sctx, "PUT", "/tag", reqBody, true, nil)
	if !isSuccess {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}
