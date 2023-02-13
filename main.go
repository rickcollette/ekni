package main

import (
	"net/http"
	"strings"

	"ekni/actions"

	"github.com/gorilla/mux"
)

func main() {
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

	http.ListenAndServe(":8675", r)

}
