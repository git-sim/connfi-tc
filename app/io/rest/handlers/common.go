package handlers

import "net/http"

// SetupCORS Cross Origin request
func SetupCORS(w http.ResponseWriter) {
	// allow all for now - probably a massive security risk
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
