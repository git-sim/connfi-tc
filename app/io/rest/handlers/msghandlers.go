package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/git-sim/tc/app/domain/entity"
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
		accIDString := session.Values["id"].(string)

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
			//Enq the message
			outmsgid, err := mu.EnqueueMsg(inmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// todo returning the ingested message in the response, could just
			//    return the id and let the client come get it.
			outmsg, err := mu.RetrieveMsg(outmsgid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(&outmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case http.MethodGet:
			r.ParseForm()
			msgIDString := r.FormValue("msgid")
			id, err := strconv.ParseUint(msgIDString, entity.AccountIDStringBase, entity.AccountIDBits)
			if err != nil {
				errStr := fmt.Sprintf("Get message invalid acc id %s, message id %s, err %s",
					accIDString, msgIDString, err)
				http.Error(w, errStr, http.StatusBadRequest)
			}
			outmsg, err := mu.RetrieveMsg(usecase.MsgIDType(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			err = json.NewEncoder(w).Encode(&outmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
