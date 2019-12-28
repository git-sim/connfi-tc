package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/git-sim/tc/app/usecase"
)

// HandleMessage handler - Allows POSTing messages to the system, and reading a message given an id
func HandleMessage(mu usecase.MsgUsecase, u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Always returns a session
		session, _ := u.GetSession().FromReq(r)
		// Could do auth here, we're interested in getting the AccountId of the user
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		idstring := session.Values["id"].(string)

		switch r.Method {
		case http.MethodPost:
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields() // catch unwanted fields

			inmsg := &usecase.IngressMsg{}
			err := d.Decode(inmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ok, err := mu.IsValid(inmsg)
			if ok && err == nil {
				//Enq the message
				outmsg, err := mu.EnqueueMsg(inmsg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = json.NewEncoder(w).Encode(&outmsg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}

		case http.MethodGet:
			id, err := strconv.ParseUint(idstring, 10, 64)
			if err != nil {
				errStr := fmt.Sprintf("Get message invalid message id %s, err %s",
					idstring, err)
				http.Error(w, errStr, http.StatusBadRequest)
			}
			outmsg, err := mu.RetrieveMsg(usecase.MsgIDType(id))
			if err != nil {
				http.Error(w, "message not found", http.StatusNotFound)
				return
			}
			err = json.NewEncoder(w).Encode(&outmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
