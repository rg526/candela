package cdsite

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
