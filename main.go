package main

import (
	"ekni/actions"
	"ekni/shared"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-ini/ini"
	"github.com/gorilla/mux"
)

func main() {
	// Load the configuration file
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	var SystemConfig shared.EkniConfig
	SystemConfig.OtpIssuer = cfg.Section("otp").Key("issuer").String()
	SystemConfig.OtpDuration, err = cfg.Section("otp").Key("duration").Int()
	if err != nil {
		log.Fatalf("Fail to parse OtpDuration: %v", err)
	}
	SystemConfig.AllowRegistration, err = cfg.Section("registration").Key("allow").Bool()
	if err != nil {
		log.Fatalf("Fail to parse AllowRegistration: %v", err)
	}
	SystemConfig.AllowRegistrationOnlyFromDomain, err = cfg.Section("registration").Key("allow_only_from_domain").Bool()
	if err != nil {
		log.Fatalf("Fail to parse AllowRegistrationOnlyFromDomain: %v", err)
	}
	SystemConfig.RegistrationDomain = cfg.Section("registration").Key("domain").String()
	SystemConfig.WireGuardPort, err = cfg.Section("wireguard").Key("port").Int()
	if err != nil {
		log.Fatalf("Fail to parse WireGuardPort: %v", err)
	}
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/web/") {
				http.ServeFile(w, r, "public"+r.URL.Path)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Endpoints
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/index.html", http.StatusFound)
	})

	r.HandleFunc("/api/init", func(w http.ResponseWriter, r *http.Request) {
		var request shared.InitRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		actions.InitWireGuardServer(w, r, request.IPAddress, fmt.Sprintf("%d", request.ListenPort), request.PrivateKey)
	}).Methods("POST")

	r.HandleFunc("/api/addclient", func(w http.ResponseWriter, r *http.Request) {
		var client shared.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		actions.CreateWireGuardClientConfig(w, r, client.Name, client.IP, client.Key)
	}).Methods("POST")

	r.HandleFunc("/api/getuserip", actions.GetUserIP).Methods("GET")

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var user shared.WebUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		actions.Login(w, r, user.Username, user.Password)
	}).Methods("POST")

	r.HandleFunc("/api/logoff/{username}", actions.Logoff).Methods("POST")
	r.HandleFunc("/api/adduser/{username}/{password}/{mfa}", func(w http.ResponseWriter, r *http.Request) { actions.AddUser(w, r, SystemConfig) }).Methods("GET")
	port := fmt.Sprintf(":%d", SystemConfig.WireGuardPort)

	// Start the HTTP server
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Fail to start HTTP server: %v", err)
	}
}
