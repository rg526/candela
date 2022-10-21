package cdsite

import (
	"io"
	"bytes"
	"errors"
	"encoding/json"
	"encoding/base64"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)


func ReportError(ctx *gin.Context, code int, title string, err error) {
	ctx.HTML(code, "layout/error", gin.H{
		"Title": "Error",
		"ErrorTitle": title,
		"ErrorDescription": "Error: " + err.Error()})
}

func ReportErrorToJSON(ctx *gin.Context, code int, title string, err error) {
	ctx.JSON(code, gin.H{
		"Status": "ERROR",
		"Error": err.Error()})
}


// Request to CDEngine
func CDRequestImpl(ctx *gin.Context, sctx *Context,
		reqType string,
		path string, value map[string]interface{}, useToken bool,
		result any,
		errCb func(*gin.Context, int, string, error)) bool {

	// Prepare request value
	reqUrl := sctx.Conf.CDAPIUrl + path
	var reqBody io.Reader

	if (reqType == "GET") {
		urlValue := url.Values{}
		for key, elem := range value {
			urlValue.Add(key, elem.(string))
		}
		reqUrl += "?" + urlValue.Encode()
	} else if (reqType == "POST" || reqType == "PUT") {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(value)
		if err != nil {
			errCb(ctx, http.StatusInternalServerError,
				"Service Error", err)
			return false
		}
		reqBody = &buf
	}

	req, err := http.NewRequest(reqType, reqUrl, reqBody)
	if err != nil {
		errCb(ctx, http.StatusInternalServerError,
			"Service Error", err)
		return false
	}

	// Add auth header?
	if (useToken) {
		session := sessions.Default(ctx)
		token := session.Get("token")
		if token == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/auth" +
				"?ret=" +
				base64.URLEncoding.EncodeToString([]byte(
					ctx.Request.URL.Path)))
			return false
		}
		req.Header.Add("Authorization", token.(string))
	}

	// Do request
	res, err := sctx.Client.Do(req)
	if err != nil {
		errCb(ctx, http.StatusBadGateway,
			"Service Error", err)
		return false
	}
	if res.StatusCode == http.StatusUnauthorized {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth" +
			"?ret=" +
			base64.URLEncoding.EncodeToString([]byte(
				ctx.Request.URL.Path)))
		return false
	}
	if res.StatusCode != http.StatusOK {
		var msg map[string]interface{}
		json.NewDecoder(res.Body).Decode(&msg)
		errCb(ctx, http.StatusBadGateway,
			"Service Error",
			errors.New("CDEngine error: " + msg["Error"].(string)))
		return false
	}

	// Decode request
	if result == nil {
		// Skip decode
		return true
	}
	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		errCb(ctx, http.StatusInternalServerError,
			"Service Error", err)
		return false
	}

	return true
}

// Request to CDEngine, Report error as HTML
func CDRequest(ctx *gin.Context, sctx *Context,
		reqType string,
		path string, value map[string]interface{}, useToken bool,
		result any) bool {

	return CDRequestImpl(ctx, sctx, reqType, path, value, useToken, result,
		ReportError);
}

// Request to CDEngine, Report error as JSON
func CDRequestErrJSON(ctx *gin.Context, sctx *Context,
		reqType string,
		path string, value map[string]interface{}, useToken bool,
		result any) bool {

	return CDRequestImpl(ctx, sctx, reqType, path, value, useToken, result,
		ReportErrorToJSON);
}
