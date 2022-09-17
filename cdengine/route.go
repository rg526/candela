package main

import (
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"time"

	"candela/cdmodel"

	"github.com/gin-gonic/gin"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func getCourse(ctx *gin.Context, db *sql.DB) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{
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


func getProfessor(ctx *gin.Context, db *sql.DB) {
	// Find course ID
	var prof cdmodel.Professor
	prof_name := ctx.Query("name")

	// Query DB
	stmtProf, err := db.Prepare("SELECT name, RMPRatingClass, RMPRatingOverall FROM professor WHERE name = ?")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		"status": "OK",
		"data": prof})
}


type config struct {
	Host	string
	Port	int
	DBUser	string
	DBPwd	string
	DBName	string
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
		getCourse(c, db)
	})
	r.GET("/professor", func(c *gin.Context) {
		getProfessor(c, db)
	})

	// Run CDENGINE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
