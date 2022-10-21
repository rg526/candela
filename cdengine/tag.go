package cdengine

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Tags

// Endpoint "/tag" PUT
// Create a new tag
// Request body:
// - CID (string)
// - Content (string)
func PutTag(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Get request body
	var reqBody struct {
		CID			string
		Content		string
	}
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ReportError(ctx, http.StatusBadRequest, err)
		return
	}

	// Do query
	_, err = ectx.DB.
		Exec("INSERT INTO tag (cid, uid, content, time) VALUES (?, ?, ?, ?)",
			reqBody.CID, user.UID, reqBody.Content, time.Now().Unix())
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK"})
}
