package actions

import (
	"ekni/shared"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Logoff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
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

	// Update the user's active status in the database
	_, err = db.Exec("UPDATE users SET active=? WHERE username=?", false, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
