package cdengine

import (
	"github.com/gin-gonic/gin"
)


func ReportError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, gin.H{
		"Status": "ERROR",
		"Error": "Error: " + err.Error()})
}

func ReportErrorFromString(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"Status": "ERROR",
		"Error": msg})
}
