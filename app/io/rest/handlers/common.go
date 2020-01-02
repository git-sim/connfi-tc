package handlers

import (
	"net/http"

	"github.com/git-sim/tc/app/usecase"
)

// SetupCORS Cross Origin request
func SetupCORS(r *http.Request, w http.ResponseWriter) {
	//
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// GetAccIDFromSession checks the session cookie returns the session info(accIDString, ok, aut)
func GetAccIDFromSession(u usecase.AccountUsecase, r *http.Request) (accIDString string, ok bool, auth bool) {

	accIDString = ""
	ok = false
	auth = false

	//Always returns a session
	session, _ := u.GetSession().FromReq(r)
	// Could do auth here, we're interested in getting the AccountId of the user
	if auth, ok = session.Values["authenticated"].(bool); ok && auth {
		accIDString, ok = session.Values["id"].(string)
	}

	// Try from url params, should be able to get rid of the above
	if ok == false {
		r.ParseForm()
		accIDString = r.FormValue("accid")
		if u.IsRegisteredID(accIDString) {
			ok = true
			auth = true
		}
	}
	return accIDString, ok, auth
}
