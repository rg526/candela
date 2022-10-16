package cdsite

import (
	"io"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

// Request to CDEngine
func CDRequest(ctx *gin.Context, sctx *Context,
		reqType string,
		path string, value map[string]string, useToken bool,
		result any) bool {

	// Prepare request value
	reqUrl := sctx.Conf.CDAPIUrl + path
	var reqBody io.Reader

	if (reqType == "GET") {
		urlValue := url.Values{}
		for key, elem := range value {
			urlValue.Add(key, elem)
		}
		reqUrl += "?" + urlValue.Encode()
	} else if (reqType == "POST") {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(value)
		if err != nil {
			ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Service Error",
				"ErrorDescription": "Connection error: " + err.Error()})
			return false
		}
		reqBody = &buf
	}

	req, err := http.NewRequest("GET", reqUrl, reqBody)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Connection error: " + err.Error()})
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
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Connection error: " + err.Error()})
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
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "CDEngine error: " + msg["Error"].(string)})
		return false
	}

	// Decode request
	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Decode error: " + err.Error()})
		return false
	}

	return true
}

func TmplTruncDesc(desc string) string {
	return desc
}
