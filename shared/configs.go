package shared

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("secret-key"))

type EkniConfig struct {
	OtpIssuer                       string
	OtpDuration                     int
	AllowRegistration               bool
	AllowRegistrationOnlyFromDomain bool
	RegistrationDomain              string
	WireGuardPort                   int
}

type WebUser struct {
	Username string
	Email    string
	Password string
	Mfa      bool
	Active   bool
	Admin    bool
}

type Client struct {
	Name string
	IP   string
	Key  string
}
