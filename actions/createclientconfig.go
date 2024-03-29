package actions

import (
	"ekni/shared"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/jmoiron/sqlx"
)

// CreateWireGuardClientConfig creates a configuration for a WireGuard client
func CreateWireGuardClientConfig(w http.ResponseWriter, r *http.Request, clientName string, clientIP string, clientPrivateKey string) {
	// Parameters are now passed as arguments, so no need to retrieve them from the URL or request body

	// Generate the WireGuard client configuration file
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s

[Peer]
PublicKey = %s
AllowedIPs = %s
`, clientPrivateKey, clientIP, clientName)
	// Save the configuration to a file
	err := os.WriteFile(fmt.Sprintf("%s.conf", clientName), []byte(config), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the client configuration to the SQLite3 database
	db, err := sqlx.Connect("sqlite3", "development.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	client := &shared.Client{
		Name: clientName,
		IP:   clientIP,
		Key:  clientPrivateKey,
	}
	_, err = db.Exec("INSERT INTO clients (name, ip, key) VALUES (?, ?, ?)", client.Name, client.IP, client.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rebuild the WireGuard server configuration
	_, err = exec.Command("wg-quick", "down", "wg0.conf").Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clients := []shared.Client{}
	err = db.Select(&clients, "SELECT * FROM clients")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverConfig := "[Interface]\n"
	for _, client := range clients {
		serverConfig += fmt.Sprintf("Peer = %s\nAllowedIPs = %s\n", client.Key, client.IP)
	}

	// Save the server configuration to a file
	err = os.WriteFile("wg0.conf", []byte(serverConfig), 0644)
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
	fmt.Fprint(w, "WireGuard client configuration created and saved successfully")
}
