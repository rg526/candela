package cdengine

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
