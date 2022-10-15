package cdengine

import (
	"log"
	"time"
	"encoding/json"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host				string
	Port				int
	DBUser				string
	DBPwd				string
	DBName				string
	OAuth2ClientID		string
	OAuth2ClientSecret	string
	OAuth2Scope			string
	OAuth2RedirectURI	string
	MaxSearchResult		int
}

type Context struct {
	DB					*sql.DB
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

	// Open DB
	ctx.DB, err = sql.Open("mysql",
		ctx.Conf.DBUser + ":" + ctx.Conf.DBPwd + "@/" +
		ctx.Conf.DBName + "?autocommit=true")
	if err != nil {
		log.Fatal("Error: open database: ", err)
	}
	ctx.DB.SetConnMaxLifetime(time.Minute * 3)

	return ctx
}
