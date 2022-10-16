package cdsite

import (
	_ "log"
	"encoding/base64"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)


// Endpoint "auth"
// Authenticate new user
func GetAuth(ctx *gin.Context, sctx *Context) {
	retPath := ctx.Query("ret")

	// Redirect to Google
	authVal := url.Values{}
	authVal.Add("client_id", sctx.Conf.OAuth2ClientID)
	authVal.Add("redirect_uri", sctx.Conf.OAuth2RedirectURI)
	authVal.Add("state", retPath)
	authVal.Add("response_type", "code")
	authVal.Add("scope", sctx.Conf.OAuth2Scope)
	authVal.Add("hd", "andrew.cmu.edu")
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth?" + authVal.Encode()

	ctx.Redirect(http.StatusMovedPermanently, authUrl)
}


// Endpoint "/authCallback"
// Callback from authentication
func GetAuthCallback(ctx *gin.Context, sctx *Context) {
	// Get authCode
	if ctx.Query("error") != "" {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Google returns " + ctx.Query("error") + "."})
		return
	}
	authCode := ctx.Query("code")

	// Verify using CDEngine
	authVal := map[string]interface{} {
		"code": authCode}
	var userResp map[string]interface{}
	isSuccess := CDRequest(ctx, sctx, "GET", "/auth", authVal, false, &userResp)
	if !isSuccess {
		return
	}

	// Set cookie
	session := sessions.Default(ctx)
	session.Set("token", userResp["Token"].(string))
	session.Save()

	// Page
	retPath := ctx.Query("state")
	data, err := base64.URLEncoding.DecodeString(retPath)
	if err != nil {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Error decoding state " + err.Error() + "."})
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, string(data))
}


// Endpoint "/logout"
// Clears current session
func GetLogout(ctx *gin.Context, sctx *Context) {
	// Remove session
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.HTML(http.StatusOK, "layout/error", gin.H{
			"Title": "Logout",
			"ErrorTitle": "You are now logged out",
			"ErrorDescription": "Please close your browswer window."})
}
