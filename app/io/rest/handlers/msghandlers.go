package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/git-sim/tc/app/usecase"
)

// HandleMessage handler - Allows POSTing messages to the system, and reading a message given an id
func HandleMessage(mu usecase.MsgUsecase, ufo usecase.FoldersUsecase, u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetupCORS(r, w)
		accIDString, ok, auth := GetAccIDFromSession(u, r)
		if !auth || !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

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
			//    return the id for brevity. The client can get the whole message if they want it.
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
			mid, err := parseIDStringAndReportErr(w, accIDString, msgIDString)
			if err != nil {
				return //error already reported
			}

			outmsg, err := mu.RetrieveMsg(usecase.MsgIDType(mid))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			err = json.NewEncoder(w).Encode(&outmsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPut:
			// Puts are control messages such as:
			//    marking a message as viewed,starred
			//    moving a message to a different folder a message
			// Fields for a put are
			//    Mark As Viewed, Starred
			//        msgid: string base16 specifies the message id
			//        viewed:    0|1
			//        starred:   0|1
			//    Move to folder
			//        msgid: same as above
			//        dest: {0|inbox, 1|archive, 2|sent, 3|scheduled}
			r.ParseForm()
			msgIDString := r.FormValue("msgid")
			mid, err := parseIDStringAndReportErr(w, accIDString, msgIDString)
			if err != nil {
				return //error already reported
			}

			accID, err := usecase.ToAccountID(accIDString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if formval, ok := r.Form["viewed"]; ok {
				newval := (formval[0] == "1")
				ufo.UpdateViewed(accID, mid, newval)
			}

			if formval, ok := r.Form["starred"]; ok {
				newval := (formval[0] == "1")
				ufo.UpdateStarred(accID, mid, newval)
			}
			// Todo handle Move to Folder Operation
		case http.MethodOptions:
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Helpers
func parseIDStringAndReportErr(w http.ResponseWriter, accIDString string, msgIDString string) (usecase.MsgIDType, error) {
	mid, err := usecase.ToMsgID(msgIDString)
	if err != nil {
		errStr := fmt.Sprintf("Get message invalid acc id %s, message id %s, err %s",
			accIDString, msgIDString, err)
		http.Error(w, errStr, http.StatusBadRequest)
	}
	return mid, err
}
