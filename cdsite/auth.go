package cdsite

import (

	"encoding/json"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

func GetAuth(ctx *gin.Context, conf Config) {
	// Redirect to Google
	authVal := url.Values{}
	authVal.Add("client_id", conf.OAuth2ClientID)
	authVal.Add("redirect_uri", conf.OAuth2RedirectURI)
	authVal.Add("response_type", "code")
	authVal.Add("scope", conf.OAuth2Scope)
	authVal.Add("hd", "andrew.cmu.edu")
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth?" + authVal.Encode()

	ctx.Redirect(http.StatusMovedPermanently, authUrl)
}

type AuthResponse struct {
	AccessToken		string		`json:"access_token"`
}

type UserResponse struct {
	Status			string
	Token			string
}

func GetAuthCallback(ctx *gin.Context, conf Config) {
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
	res, err := http.Get(conf.CDAPIUrl + "auth?" + authVal.Encode())
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
	var user UserResponse
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		ctx.HTML(http.StatusOK, "layout/error", gin.H{
				"Title": "Error",
				"ErrorTitle": "Auth Error",
				"ErrorDescription": "Error decoding json" + err.Error() + "."})
		return
	}

	// Set cookie
	session := sessions.Default(ctx)
	session.Set("token", user.Token)
	session.Save()

	// Page
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func GetLogout(ctx *gin.Context, conf Config) {
	// Remove session
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.HTML(http.StatusOK, "layout/error", gin.H{
			"Title": "Logout",
			"ErrorTitle": "You are now logged out",
			"ErrorDescription": "Please close your browswer window."})
}
