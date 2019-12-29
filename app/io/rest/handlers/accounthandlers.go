package handlers

import (
	"fmt"
	"net/http"

	"github.com/git-sim/tc/app/usecase"
)

func HandleAccount(u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Always returns a session
		session, _ := u.GetSession().FromReq(r)
		// Could do auth here, we're interested in getting the AccountId of the user
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		accIDString := session.Values["id"].(string)

		r.ParseForm()
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "missing email in request", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPost:
			acc, err := u.RegisterAccount(email)
			if err != nil {
				http.Error(w, "email already in use", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintf(w, "id: %s\n", acc.ID)
		case http.MethodDelete:
			err := u.DeleteAccount(email)
			if err != nil {

				http.Error(w, "email not found", http.StatusNotFound)
				return
			}

		case http.MethodGet:
			acc, err := u.GetAccount(email)
			if err != nil {
				http.Error(w, "email not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "id: %s, email: %s, FirstName: %s, LastName: %s\n",
				acc.ID, acc.Email, acc.FirstName, acc.LastName)
			//w.Write(acc)

		case http.MethodPut:
			acc, err := u.GetAccount(email)
			if err != nil {
				http.Error(w, "email not found", http.StatusNotFound)
				return
			}
			// Verify that Check that the logged in account is the one that's making the change
			if accIDString == acc.ID {
				// Have to distinguish between an empty name and a non-specified name
				// Don't change the name if it was nil
				var pfirstname *string
				var plastname *string
				fnval, ok := r.Form["firstname"]
				if ok {
					pfirstname = &fnval[0]
				}

				lnval, ok := r.Form["lastname"]
				if ok {
					plastname = &lnval[0]
				}

				err := u.UpdateNameAccount(email, pfirstname, plastname)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}

func HandleAccountList(u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			accs, err := u.GetAccountList()
			if err != nil {
				http.Error(w, "email not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "count: %d\n", len(accs))
			for _, acc := range accs {
				_, _ = fmt.Fprintf(w, "id: %s, email: %s\n", acc.ID, acc.Email)
			}
			//w.Write(acc)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
