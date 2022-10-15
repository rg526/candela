package cdsite

import (
	"encoding/json"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	"candela/cdmodel"
)


func VerifyUserFromSession(ctx *gin.Context, sctx *Context) (string, cdmodel.User, bool) {
	// Get token from session
	session := sessions.Default(ctx)
	token := session.Get("token")
	if token == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth")
		return "", cdmodel.User{}, false
	}

	// Find user by token
	userVal := url.Values{}
	userVal.Add("token", token.(string))
	userUrl := sctx.Conf.CDAPIUrl + "user?" + userVal.Encode()

	// Send CDAPI request
	var userResp struct {
		Status		string
		Data		cdmodel.User
	}
	res, err := http.Get(userUrl)
	if err != nil {
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Unsuccessful connection to CDEngine."})
		return "", cdmodel.User{}, false
	}
	if res.StatusCode == http.StatusUnauthorized {
		ctx.Redirect(http.StatusTemporaryRedirect, "/auth")
		return "", cdmodel.User{}, false
	}
	if res.StatusCode != http.StatusOK {
		var msg map[string]interface{}
		json.NewDecoder(res.Body).Decode(&msg)
		ctx.HTML(http.StatusBadGateway, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "CDEngine error: " + msg["Error"].(string)})
		return "", cdmodel.User{}, false
	}
	err = json.NewDecoder(res.Body).Decode(&userResp)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "layout/error", gin.H{
			"Title": "Error",
			"ErrorTitle": "Service Error",
			"ErrorDescription": "Decode error: " + err.Error()})
		return "", cdmodel.User{}, false
	}

	// Return user struct
	return token.(string), userResp.Data, true
}


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
