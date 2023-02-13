package actions

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type WebUser struct {
	Username string
	Email    string
	Password string
	Mfa      bool
	Active   bool
	Admin    bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	token := r.FormValue("token")

	// Authenticate the user's password
	db, err := sqlx.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user := WebUser{}
	err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		http.Error(w, "Username or password incorrect", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Username or password incorrect", http.StatusUnauthorized)
		return
	}

	if !user.Active {
		http.Error(w, "Your account has been deactivated", http.StatusUnauthorized)
		return
	}

	if !user.Mfa {
		// User does not have MFA enabled, proceed to login
		// ...
		return
	}

	// User has MFA enabled, verify the token
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example Inc.",
		AccountName: username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valid, err := totp.ValidateCustom(token, secret.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Token is valid, proceed to login
	// ...
}
