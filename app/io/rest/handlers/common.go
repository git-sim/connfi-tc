package handlers

import "net/http"

// SetupCORS Cross Origin request
func SetupCORS(r *http.Request, w http.ResponseWriter) {
	//
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
}
