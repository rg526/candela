package cdengine

import (
	"database/sql"
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
}

type Context struct {
	DB					*sql.DB
	Conf				Config
}
