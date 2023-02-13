package actions

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

// InitWireGuardServer initializes a WireGuard server with the given parameters
func InitWireGuardServer(w http.ResponseWriter, r *http.Request) {
	// Retrieve the configuration parameters from the request
	vars := mux.Vars(r)
	ipAddress := vars["ip_address"]
	listenPort := vars["listen_port"]
	privateKey := vars["private_key"]
	// Generate the WireGuard server configuration file
	config := fmt.Sprintf(`[Interface]
ddress = %s
ListenPort = %s
PrivateKey = %s
`, ipAddress, listenPort, privateKey)
	// Save the configuration to a file
	err := ioutil.WriteFile("wg0.conf", []byte(config), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Bring up the WireGuard interface using the wg-quick tool
	_, err = exec.Command("wg-quick", "up", "wg0.conf").Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("WireGuard server configuration initialized successfully"))
	db, err := sql.Open("sqlite3", "./migrations.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create the clients table
	createClientsTable := `
	CREATE TABLE IF NOT EXISTS clients (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name VARCHAR(255) NOT NULL,
	  ip VARCHAR(255) NOT NULL,
	  key VARCHAR(255) NOT NULL
	);
	`
	_, err = db.Exec(createClientsTable)
	if err != nil {
		log.Fatalf("Error creating clients table: %v", err)
	}

	// Create the WebUser table
	createWebUserTable := `
	CREATE TABLE IF NOT EXISTS web_users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    username VARCHAR(255) NOT NULL,
	    email VARCHAR(255) NOT NULL,
	    password VARCHAR(255) NOT NULL,
	    mfa BOOLEAN NOT NULL,
	    active BOOLEAN NOT NULL,
	    admin BOOLEAN NOT NULL
	);
	`
	_, err = db.Exec(createWebUserTable)
	if err != nil {
		log.Fatalf("Error creating web_users table: %v", err)
	}

	// Create the web_user_clients table
	createWebUserClientsTable := `
	CREATE TABLE IF NOT EXISTS web_user_clients (
	  web_user_id INTEGER NOT NULL,
	  client_id INTEGER NOT NULL,
	  PRIMARY KEY (web_user_id, client_id),
	  FOREIGN KEY (web_user_id) REFERENCES web_users (id),
	  FOREIGN KEY (client_id) REFERENCES clients (id)
	);
	`
	_, err = db.Exec(createWebUserClientsTable)
	if err != nil {
		log.Fatalf("Error creating web_user_clients table: %v", err)
	}

	log.Println("Migrations complete")
}
