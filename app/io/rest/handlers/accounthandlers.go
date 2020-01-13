package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/git-sim/tc/app/usecase"
)

// AccountCtxFunc returns a context handler to validate the accountID exists, and is registered
func AccountCtxFunc(u usecase.AccountUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlAccountID := chi.URLParam(r, "accountID")
			ok := u.IsRegisteredID(urlAccountID)
			if !ok {
				fmt.Printf("AccountCtxFunc Account Not found %s\n", urlAccountID)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			// Save away the accountID for this request so later handlers don't have
			// to keep validating it.  Applies this request only
			accountID, err := usecase.ToAccountID(urlAccountID)
			if err != nil {
				ReportUsecaseFault(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "dbAccountID", accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// POST /accounts
func CreateAccount(u usecase.AccountUsecase) http.HandlerFunc {
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
		acc, err := u.GetAccount(accountID)
		if err != nil {
			ReportUsecaseFault(w, err)
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

// PUT  /accounts/1234
func PutAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
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
			ReportUsecaseFault(w, err)
			return
		}
		return
	}
}

// DEL  /accounts/1234
func DeleteAccount(u usecase.AccountUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		err := u.DeleteAccount(accountID)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}
		return
	}
}
