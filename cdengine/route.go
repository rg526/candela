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
		log.Fatal("Error", err)
	}

	// Query DB
	stmtCourse, err := db.Prepare("SELECT cid, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount FROM course WHERE cid = ?")
	if err != nil {
		log.Fatal("Error", err)
	}
	err = stmtCourse.QueryRow(cid).Scan(&course.CID, &course.Description, &course.Dept, &course.Units, &course.Prof, &course.Prereq, &course.Coreq, &course.FCEHours, &course.FCETeachingRate, &course.FCECourseRate, &course.FCELevel, &course.FCEStudentCount)
	if err != nil {
		log.Fatal("Error", err)
	}

	// Return result
	ctx.JSON(http.StatusOK, course)
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

	// Run CDENGINE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
