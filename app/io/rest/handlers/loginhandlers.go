package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/git-sim/tc/app/usecase"
)

func exampleSessionID(us usecase.SessionUsecase, r *http.Request, w http.ResponseWriter) {
	session, _ := us.FromReq(r)
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		//http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	id := session.Values["id"]
	w.Write(id.([]byte))
	return
}

// HandleLogin - handles logging in or registering a new account
func HandleLogin(us usecase.SessionUsecase, u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetupCORS(r,w)
		switch r.Method {
		case http.MethodPost:
			r.ParseForm()
			email := r.FormValue("email")
			if email == "" {
				http.Error(w, "missing email in request", http.StatusBadRequest)
				return
			}

			session, _ := us.FromReq(r)

			acc, err := u.GetAccount(email)
			if err != nil {
				es, ok := err.(*usecase.ErrStat)
				if ok && es.Code == usecase.EsNotFound {
					// Doesn't exist create new account
					acc, err = u.RegisterAccount(email)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						fmt.Fprintf(w, "err: %s\n", err.Error())
						return
					}
					//w.WriteHeader(http.StatusCreated) //the cookie is set iff StatusOk
				} else {
					http.Error(w, "Error while retrieving account", http.StatusInternalServerError)
					fmt.Fprintf(w, "err: %s\n", err.Error())
					return
				}
			}
			if acc != nil {

				// Any real authentication would potentially go here
				session.Values["authenticated"] = true
				session.Values["id"] = acc.ID
				session.Save(r, w)

				err = json.NewEncoder(w).Encode(acc)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}

// HandleLogout clears out the session id
func HandleLogout(us usecase.SessionUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetupCORS(r,w)
		r.ParseForm()
		session, _ := us.FromReq(r)

		session.Values["authenticated"] = false
		session.Values["id"] = ""
		session.Save(r, w)
	})
}
