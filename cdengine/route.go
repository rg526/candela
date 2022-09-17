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

func getCourse(c *gin.Context) {
	var x cdmodel.Course
	c.JSON(http.StatusOK, x)
}


type config struct {
	Host	string
	Port	int
}

func main() {
	conf_content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	var conf config
	err = json.Unmarshal(conf_content, &conf)
	if err != nil {
		log.Fatal("Error: read config file: ", err)
	}

	r := gin.Default()
	r.GET("/course", getCourse)

	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))
}
