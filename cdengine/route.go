package main

import (
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"time"
	"net/url"

	"candela/cdmodel"

	"github.com/gin-gonic/gin"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func verifyToken(token string, db *sql.DB) (string, string, error) {
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

func getCourse(ctx *gin.Context, db *sql.DB, conf config) {
	// Find course ID
	var course cdmodel.Course
	cid_query := ctx.Query("cid")
	cid, err := strconv.Atoi(cid_query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Query DB
	stmtCourse, err := db.Prepare("SELECT cid, name, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount FROM course WHERE cid = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	err = stmtCourse.QueryRow(cid).Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq, &course.FCEHours, &course.FCETeachingRate, &course.FCECourseRate, &course.FCELevel, &course.FCEStudentCount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": course})
}


func getProfessor(ctx *gin.Context, db *sql.DB, conf config) {
	// Find course ID
	var prof cdmodel.Professor
	prof_name := ctx.Query("name")

	// Query DB
	stmtProf, err := db.Prepare("SELECT name, RMPRatingClass, RMPRatingOverall FROM professor WHERE name = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}
	err = stmtProf.QueryRow(prof_name).Scan(&prof.Name, &prof.RMPRatingClass, &prof.RMPRatingOverall)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
		return
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": prof})
}


type AuthResponse struct {
    AccessToken     string      `json:"access_token"`
}

type UserInfoResponse struct {
    UID             string      `json:"id"`
    Name            string      `json:"name"`
    GivenName       string      `json:"given_name"`
    FamilyName      string      `json:"family_name"`
    Email           string      `json:"email"`
}

func getAuth(ctx *gin.Context, db *sql.DB, conf config) {
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

	var result AuthResponse
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
        return
    }

    // Get email and name
    userVal := url.Values{}
    userVal.Add("access_token", result.AccessToken)
    res, err = http.Get("https://www.googleapis.com/oauth2/v1/userinfo?" + userVal.Encode())
    if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
        return
    }
    var userResp UserInfoResponse
    err = json.NewDecoder(res.Body).Decode(&userResp)
    if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"Status": "ERROR",
			"Error": "Error: " + err.Error()})
        return
    }

	// Generate user structure
	var user cdmodel.User
	user.UID = userResp.UID
	user.Name = userResp.Name
	user.GivenName = userResp.GivenName
	user.FamilyName = userResp.FamilyName
	user.Email = userResp.Email

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


type config struct {
	Host				string
	Port				int
	DBUser				string
	DBPwd				string
	DBName				string
	OAuth2ClientID		string
	OAuth2ClientSecret	string
	OAuth2Scope			string
	OAuth2RedirectURI	string
}

func main() {
	// Read config file
	conf_content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	var conf config
	err = json.Unmarshal(conf_content, &conf)
	if err != nil {
		log.Fatal("Error: read config file: ", err)
	}

	// Open DB
	db, err := sql.Open("mysql",
		conf.DBUser + ":" + conf.DBPwd + "@/" +
		conf.DBName + "?autocommit=true")
	if err != nil {
		log.Fatal("Error: open database: ", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)


	// Setup routes
	r := gin.Default()
	r.GET("/course", func(c *gin.Context) {
		getCourse(c, db, conf)
	})
	r.GET("/professor", func(c *gin.Context) {
		getProfessor(c, db, conf)
	})
	r.GET("/auth", func(c *gin.Context) {
		getAuth(c, db, conf)
	})

	// Run CDENGINE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
