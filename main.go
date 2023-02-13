package main

import (
	"ekni/actions"
	"ekni/shared"
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
	var config shared.EkniConfig
	config.OtpIssuer = cfg.Section("otp").Key("issuer").String()
	config.OtpDuration, err = cfg.Section("otp").Key("duration").Int()
	if err != nil {
		log.Fatalf("Fail to parse OtpDuration: %v", err)
	}
	config.AllowRegistration, err = cfg.Section("registration").Key("allow").Bool()
	if err != nil {
		log.Fatalf("Fail to parse AllowRegistration: %v", err)
	}
	config.AllowRegistrationOnlyFromDomain, err = cfg.Section("registration").Key("allow_only_from_domain").Bool()
	if err != nil {
		log.Fatalf("Fail to parse AllowRegistrationOnlyFromDomain: %v", err)
	}
	config.RegistrationDomain = cfg.Section("registration").Key("domain").String()
	config.WireGuardPort, err = cfg.Section("wireguard").Key("port").Int()
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

	r.HandleFunc("/api/init/{ip_address}/{listen_port}/{private_key}", actions.InitWireGuardServer).Methods("GET")
	r.HandleFunc("/api/addclient/{client_name}/{client_ip}/{client_private_key}", actions.CreateWireGuardClientConfig).Methods("GET")
	r.HandleFunc("/api/getuserip", actions.GetUserIP).Methods("GET")
	r.HandleFunc("/api/login/{username}/{password}", actions.Login).Methods("GET")
	r.HandleFunc("/api/logoff/{username}", actions.Logoff).Methods("POST")
	r.HandleFunc("/api/adduser/{username}/{password}/{mfa}", actions.AddUser).Methods("GET")

	port := fmt.Sprintf(":%d", config.WireGuardPort)

	// Start the HTTP server
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Fail to start HTTP server: %v", err)
	}
}
