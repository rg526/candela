package main

import (
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"

	"candela/cdmodel"

	"github.com/gin-gonic/gin"
	"net/http"
)

func getCourse(ctx *gin.Context) {
	// Find course ID
	var course cdmodel.Course
	cid_query := ctx.Query("cid")
	cid, err := strconv.Atoi(cid_query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": "Error: " + err.Error()})
		return
	}
	_ = course
	_ = cid


	// TODO generate HTML
}


type config struct {
	Host	string
	Port	int
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

	// Setup routes
	r := gin.Default()
	r.GET("/course", getCourse)

	// Run CDSITE
	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}