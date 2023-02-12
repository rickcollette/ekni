package actions

import (
	"fmt"
	"io/ioutil"
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

}
