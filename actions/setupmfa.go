package actions

import (
	"ekni/shared"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pquerna/otp/totp"
)

func SetupMfa(w http.ResponseWriter, r *http.Request, SystemConfig shared.EkniConfig) {
	username := r.FormValue("username")
	// Authenticate the user
	db, err := sqlx.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user := shared.WebUser{}
	err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Mfa {
		http.Error(w, "MFA is already enabled for this user", http.StatusBadRequest)
		return
	}

	// Generate a TOTP secret for the user
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      SystemConfig.OtpIssuer,
		AccountName: username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the TOTP secret to the database
	_, err = db.Exec("UPDATE users SET mfa=?, secret=? WHERE username=?", true, secret.Secret(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the QR code URL to the client
	url := secret.URL()
	w.Write([]byte(url))
}
