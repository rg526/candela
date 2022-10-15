package cdsite

import (

	"encoding/json"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)


func GetAuth(ctx *gin.Context, sctx *Context) {
	// Redirect to Google
	authVal := url.Values{}
	authVal.Add("client_id", sctx.Conf.OAuth2ClientID)
	authVal.Add("redirect_uri", sctx.Conf.OAuth2RedirectURI)
	authVal.Add("response_type", "code")
	authVal.Add("scope", sctx.Conf.OAuth2Scope)
	authVal.Add("hd", "andrew.cmu.edu")
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth?" + authVal.Encode()

	ctx.Redirect(http.StatusMovedPermanently, authUrl)
}


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
	authVal := url.Values{}
	authVal.Add("code", authCode)
	res, err := http.Get(sctx.Conf.CDAPIUrl + "auth?" + authVal.Encode())
	if err != nil {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Service Error",
				"ErrorDescription": "Unsuccessful connection to CDEngine."})
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Invalid request."})
		return
	}
	var userResp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&userResp)
	if err != nil {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Error decoding json" + err.Error() + "."})
		return
	}

	// Set cookie
	session := sessions.Default(ctx)
	session.Set("token", userResp["Token"].(string))
	session.Save()

	// Page
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

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
