package cdengine

import (
	"encoding/json"
	"time"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	
	"candela/cdmodel"
)

func VerifyToken(token string, db *sql.DB) (string, string, error) {
	// Check if token exists
	// Query DB
	stmtToken, err := db.Prepare("SELECT uid, time FROM token WHERE token = ?")
	if err != nil {
		return "", "", err
	}

	// Record UID and time
	var uid, time string
	err = stmtToken.QueryRow(token).Scan(&uid, &time)
	if err != nil {
		return "", "", err
	}

	// Return result
	return uid, time, nil
}


func GetAuth(ctx *gin.Context, db *sql.DB, conf Config) {
	// Verify code
	authCode := ctx.Query("code")

	// Turn to OAuth token
    authVal := url.Values{}
    authVal.Add("client_id", conf.OAuth2ClientID)
    authVal.Add("client_secret", conf.OAuth2ClientSecret)
    authVal.Add("code", authCode)
    authVal.Add("grant_type", "authorization_code")
    authVal.Add("redirect_uri", conf.OAuth2RedirectURI)
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

	// Test if user exists
	stmtCheck, err := db.Prepare("SELECT uid FROM user WHERE uid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	rows, err := stmtCheck.Query(user.UID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	if !rows.Next() {
		// Insert user
		stmtInsert, err := db.Prepare("INSERT INTO user (uid, name, givenName, familyName, Email) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Status": "ERROR",
				"Error": "Error: " + err.Error()})
			return
		}
		_, err = stmtInsert.Exec(user.UID, user.Name, user.GivenName ,user.FamilyName, user.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Status": "ERROR",
				"Error": "Error: " + err.Error()})
			return
		}
	}

	// Insert token
	stmtInsert, err := db.Prepare("INSERT INTO token (token, uid, time) VALUES (?, ?, ?)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	_, err = stmtInsert.Exec(userToken, user.UID, time.Now().Unix())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return token
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Token": userToken})
}
