package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/git-sim/tc/app/usecase"
)

// POST /accounts
func CreateAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST /accounts")

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

		if account.Email == "" {
			http.Error(w, "Email not specified", http.StatusBadRequest)
			return
		}

		acc, err := u.RegisterAccount(&account)
		if err != nil {
			if usecase.CheckEs(err, usecase.EsAlreadyExists) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
		return
	}
}

// GET  /accounts
func GetAccountList(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /accounts")
		accs, err := u.GetAccountList()
		if err != nil {
			http.Error(w, "account list not found", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(accs)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
		return
	}
}

// GET  /accounts/1234
func GetAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		fmt.Printf("GET /accounts/%s", accountID)
		ok := u.IsRegisteredID(accountID)
		if !ok {
			http.Error(w, "account not registered", http.StatusNotFound)
			return
		}

		acc, err := u.GetAccount(accountID)
		if err != nil {
			var es *usecase.ErrStat
			if errors.As(err, &es) {
				http.Error(w, es.Msg, ErrStatToHTTPCode(es))
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}

		return
	}
}

// PUT  /accounts/1234
func PutAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		fmt.Printf("PUT /accounts/%s", accountID)

		ok := u.IsRegisteredID(accountID)
		if !ok {
			http.Error(w, "account not registered", http.StatusNotFound)
			return
		}

		var account usecase.Account
		err := decodeJSONBody(w, r, &account)
		if err != nil {
			var es *ErrorJSONDecode
			if errors.As(err, &es) {
				http.Error(w, es.msg, es.status)
			} else {
				ReportJSONFault(w, err)
			}
			return
		}

		//fill in the body from the URL since that's the authoritative id
		account.ID = accountID
		err = u.UpdateAccount(&account)
		if err != nil {
			var es *usecase.ErrStat
			if errors.As(err, &es) {
				http.Error(w, es.Msg, ErrStatToHTTPCode(es))
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		return
	}
}

// DEL  /accounts/1234
func DeleteAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		fmt.Printf("DELETE /accounts/%s", accountID)
		ok := u.IsRegisteredID(accountID)
		if !ok {
			http.Error(w, "account not registered", http.StatusNotFound)
			return
		}

		err := u.DeleteAccount(accountID)
		if err != nil {
			var es *usecase.ErrStat
			if errors.As(err, &es) {
				http.Error(w, es.Msg, ErrStatToHTTPCode(es))
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		return
	}
}

func HandleAccount(u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		accIDString, ok, auth := getAccIDFromSession(u, r)
		if !auth || !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		r.ParseForm()
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "missing email in request", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPost:
			acc, err := u.RegisterAccountByEmail(email)
			if err != nil {
				if usecase.CheckEs(err, usecase.EsAlreadyExists) {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			err = json.NewEncoder(w).Encode(acc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			err := u.DeleteAccountByEmail(email)
			if err != nil {

				http.Error(w, "email not found", http.StatusNotFound)
				return
			}

		case http.MethodGet:
			acc, err := u.GetAccountByEmail(email)
			if err != nil {
				http.Error(w, "email not found", http.StatusNotFound)
				return
			}
			err = json.NewEncoder(w).Encode(acc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case http.MethodPut:
			acc, err := u.GetAccountByEmail(email)
			if err != nil {
				http.Error(w, "email not found", http.StatusNotFound)
				return
			}
			// Check that the logged in account is the one that's making the change
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
		setupCORS(w, r)
		switch r.Method {
		case http.MethodGet:
			accs, err := u.GetAccountList()
			if err != nil {
				http.Error(w, "account list not found", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(accs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
