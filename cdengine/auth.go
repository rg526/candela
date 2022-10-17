package cdengine

import (
	"encoding/json"
	"time"
	"net/url"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"candela/cdmodel"
)


// Verify if a token is valid
func VerifyToken(token string, ectx *Context) (cdmodel.User, error) {
	var user cdmodel.User
	// Check if token exists
	// Query DB
	// Record UID and time
	var uid, time string
	err := ectx.DB.
		QueryRow("SELECT uid, time FROM token WHERE token = ?",
			token).
		Scan(&uid, &time)
	if err != nil {
		return user, err
	}

	// Query DB for user
	err = ectx.DB.
		QueryRow("SELECT uid, name, givenName, familyName, Email FROM user WHERE UID = ?",
			uid).
		Scan(&user.UID, &user.Name, &user.GivenName, &user.FamilyName, &user.Email)
	if err != nil {
		return user, err
	}

	// Return result
	return user, nil
}


// Verify if a token is valid, from context
func VerifyTokenFromCtx(ctx *gin.Context, ectx *Context) (cdmodel.User, bool) {
	tokenArr := ctx.Request.Header["Authorization"]
	if len(tokenArr) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Status": "ERROR",
			"Error": "Error: No token in request header"})
		return cdmodel.User{}, false
	}
	token := tokenArr[0]
	user, err := VerifyToken(token, ectx)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return cdmodel.User{}, false
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return cdmodel.User{}, false
	}
	return user, true
}


// Endpoint "/user"
// Get current user info
func GetUser(ctx *gin.Context, ectx *Context) {
	// Verify token
	user, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": user})
}


// Endpoint "/auth"
// Authenticate with a OAuth code
func GetAuth(ctx *gin.Context, ectx *Context) {
	// Verify code
	authCode := ctx.Query("code")

	// Turn to OAuth token
	authVal := url.Values{}
	authVal.Add("client_id", ectx.Conf.OAuth2ClientID)
	authVal.Add("client_secret", ectx.Conf.OAuth2ClientSecret)
	authVal.Add("code", authCode)
	authVal.Add("grant_type", "authorization_code")
	authVal.Add("redirect_uri", ectx.Conf.OAuth2RedirectURI)
	res, err := http.PostForm("https://oauth2.googleapis.com/token", authVal)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + res.Status})
		return
	}

	// Decode token
	var authResp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&authResp)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Get email and name
	userVal := url.Values{}
	userVal.Add("access_token", authResp["access_token"].(string))
	res, err = http.Get("https://www.googleapis.com/oauth2/v1/userinfo?" + userVal.Encode())
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + res.Status})
		return
	}

	// Decode email and name
	var userResp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&userResp)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Generate user structure
	var user cdmodel.User
	user.UID = userResp["id"].(string)
	user.Name = userResp["name"].(string)
	user.GivenName = userResp["given_name"].(string)
	user.FamilyName = userResp["family_name"].(string)
	user.Email = userResp["email"].(string)

	// Generate token
	userToken := uuid.New().String()

	// Insert user
	_, err = ectx.DB.
		Exec("INSERT IGNORE INTO user (uid, name, givenName, familyName, Email) VALUES (?, ?, ?, ?, ?)",
			user.UID, user.Name, user.GivenName ,user.FamilyName, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Insert token
	_, err = ectx.DB.
		Exec("INSERT INTO token (token, uid, time) VALUES (?, ?, ?)",
			userToken, user.UID, time.Now().Unix())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return token
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Token": userToken})
}
