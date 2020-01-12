package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/git-sim/tc/app/usecase"
)

// Folders are the lists of messages like inbox, archive, sent, etc
// The handlers for this enpoint are about getting users requests for messages, and sending
// light control like marking a message as read/starred, or archiving message.
// The folder presenter handles the sorting, and presenting the lists of messages.

// This implementation doesn't allow adding or deleting folders.

// HandleFolder handler
func HandleFolder(ufo usecase.FoldersUsecase, mu usecase.MsgUsecase, u usecase.AccountUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		accIDString, ok, auth := getAccIDFromSession(u, r)
		if !auth || !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		accID, err := usecase.ToAccountID(accIDString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Get Folder has the params:
			//   folderid: 0|Inbox,1|Archive,2|Sent,3|Scheduled Def=0
			//   sort: 0|time,1|subject,2|sender Def=0
			//   sortorder: -1,1 Def=1
			//   limit: 0.. Def=10
			//   page: 0.. Def=0
			r.ParseForm()
			qparams := getFolderParams(r)
			pOut, err := ufo.QueryMsgs(accID, qparams)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			err = json.NewEncoder(w).Encode(pOut)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func parseIntField(in string, def, min, max int) int {
	ret, err := strconv.Atoi(in)
	if err != nil || ret < min || max < ret {
		ret = def
	}
	return ret
}

func parseSortOrder(in string) int {
	if in == "-1" {
		return 0
	}
	return 1
}

// getFolderParams parses the form intpu and gets the query parameters, defaulting out
// those that are missing or invalid.  NOTE r.ParseForm() must be called before this for it to be valid
func getFolderParams(r *http.Request) usecase.QueryParams {
	qp := usecase.QueryParams{}
	qp.FolderIdx = parseIntField(r.FormValue("folderid"), usecase.EnumInbox,
		usecase.EnumInbox, usecase.EnumNumFolders-1)
	qp.SortBy = parseIntField(r.FormValue("sort"), usecase.EnumSortByTime,
		usecase.EnumSortByTime, usecase.EnumNumSortBy-1)
	qp.SortOrder = parseSortOrder(r.FormValue("sortorder"))
	qp.Limit = parseIntField(r.FormValue("limit"), 10, 0, 100)
	qp.Page = parseIntField(r.FormValue("page"), 0, 0, 1e3)
	return qp
}

// MissingRequiredFields Returns fields that were missing
func MissingRequiredFields(r *http.Request, fields []string) []string {
	missing := make([]string, 0, len(fields))
	for _, field := range fields {
		if _, ok := r.Form[field]; !ok {
			missing = append(missing, field)
		}
	}
	return missing
}
