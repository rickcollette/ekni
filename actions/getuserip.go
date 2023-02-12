package actions

import (
	"net/http"
)

// GetUserIP returns the IP address of the user accessing the REST API
func GetUserIP(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	w.Write([]byte(ip))
}
