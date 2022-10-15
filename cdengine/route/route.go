package main

import (
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"time"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"candela/cdengine"
)


func main() {
	// Read config file
	conf_content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	var conf cdengine.Config
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
		cdengine.GetCourse(c, db, conf)
	})
	r.GET("/professor", func(c *gin.Context) {
		cdengine.GetProfessor(c, db, conf)
	})
	r.GET("/auth", func(c *gin.Context) {
		cdengine.GetAuth(c, db, conf)
	})

	// Run CDENGINE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
