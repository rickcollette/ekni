package actions

import (
	"ekni/shared"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, username string, password string) {
	// Removed the lines that extracted username and password from the request

	token := r.FormValue("token") 
	db, err := sqlx.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user := shared.WebUser{}
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
		session, err := shared.Store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["username"] = username
		session.Values["isLoggedIn"] = true
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	session, err := shared.Store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["username"] = username
	session.Values["isLoggedIn"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
