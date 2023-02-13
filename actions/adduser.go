package actions

import (
	"ekni/shared"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	mfa := r.FormValue("mfa") == "true"
	db, err := sqlx.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Check if the user already exists in the database
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM WebUser WHERE username=?", username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO WebUser (username, email, password, mfa, active, admin) VALUES (?, ?, ?, ?, ?, ?)", username, email, string(hashedPassword), mfa, true, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If mfa is true, set up MFA for the user
	if mfa {
		url, err := AddNewMfa(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(url))
	}
	w.WriteHeader(http.StatusCreated)
}

func AddNewMfa(username string) (string, error) {
	// Authenticate the user
	db, err := sqlx.Open("sqlite3", "users.db")
	if err != nil {
		return "", err
	}
	defer db.Close()
	user := shared.WebUser{}
	err = db.Get(&user, "SELECT * FROM WebUser WHERE username=?", username)
	if err != nil {
		return "", err
	}

	// Generate a TOTP secret for the user
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example Inc.",
		AccountName: username,
	})
	if err != nil {
		return "", err
	}

	// Save the TOTP secret to the database
	_, err = db.Exec("UPDATE WebUser SET mfa=?, secret=? WHERE username=?", true, secret.Secret(), username)
	if err != nil {
		return "", err
	}

	// Return the QR code URL to the client
	url := secret.URL()
	return url, nil
}
