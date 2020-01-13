package handlers

import (
	"encoding/json"
	"errors"
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

func PostLogin(us usecase.SessionUsecase, u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The account is a post body param encoded as JSON
		// Need to ensure the email is unique
		// ignore the id since a new one will be assigned
		var account usecase.Account
		err := decodeJSONBody(w, r, &account)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}

		email := account.Email
		if email == "" {
			http.Error(w, "Email not specified", http.StatusBadRequest)
			return
		}

		r.ParseForm()
		session, _ := us.FromReq(r)

		acc, err := u.GetAccountByEmail(email)
		if err != nil {
			es, ok := err.(*usecase.ErrStat)
			if ok && es.Code == usecase.EsNotFound {
				// Doesn't exist create new account
				acc, err = u.RegisterAccountByEmail(email)
				if err != nil {
					ReportUsecaseFault(w, err)
					return
				}
				//w.WriteHeader(http.StatusCreated) //the cookie is set iff StatusOk
			} else {
				ReportUsecaseFault(w, err)
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
	}
}

// PostLogout clears out the session id
func PostLogout(us usecase.SessionUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		session, _ := us.FromReq(r)

		session.Values["authenticated"] = false
		session.Values["id"] = ""
		session.Save(r, w)
	}
}
