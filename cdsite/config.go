package cdsite

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Host				string
	Port				int
	CDAPIUrl			string
	CookieSecret		string
	OAuth2ClientID		string
	OAuth2ClientSecret	string
	OAuth2Scope			string
	OAuth2RedirectURI	string
}

type Context struct {
	Conf				Config
}

func InitContext(confPath string) Context {
	var ctx Context

	// Read config file
	confContent, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("Error: open config file: ", err)
	}
	err = json.Unmarshal(confContent, &ctx.Conf)
	if err != nil {
		log.Fatal("Error: read config file: ", err)
	}

	return ctx
}
