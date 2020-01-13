package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/git-sim/tc/app/usecase"
)

// Folders are the lists of messages like inbox, archive, sent, etc
// The handlers for this enpoint are about getting users requests for messages, and sending
// light control like marking a message as read/starred, or archiving message.
// The folder presenter handles the sorting, and presenting the lists of messages.

// This implementation doesn't allow adding or deleting folders.

// FolderCtx returns a context handler to validate the folderID is valid
func FolderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlFolderID := chi.URLParam(r, "folderID")
		_, ok := checkBoundsInt(urlFolderID, usecase.EnumInbox, usecase.EnumNumFolders-1)
		if !ok {
			fmt.Printf("FolderCtxFunc Invalid folderID %s\n", urlFolderID)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		next.ServeHTTP(w, r)
	})
}

// GetFolderList returns folder info. GET /accounts/{accountID}/folders
func GetFolderList(ufo usecase.FoldersUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlAccountID := chi.URLParam(r, "accountID")
		accountID, err := usecase.ToAccountID(urlAccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		output, err := ufo.GetFolderInfo(accountID)
		if err != nil {
			http.Error(w, "Folder info not found", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
		return
	}
}

// GetFolder returns the contents of a folder with give query params GET /accounts/1234/folders/{folderID}
func GetFolder(ufo usecase.FoldersUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlAccountID := chi.URLParam(r, "accountID")
		accountID, err := usecase.ToAccountID(urlAccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get Folder has the params:
		//   sort: 0|time,1|subject,2|sender Def=0
		//   sortorder: -1,1 Def=1
		//   limit: 0.. Def=10
		//   page: 0.. Def=0
		r.ParseForm()
		qparams := getFolderParams(r)
		pOut, err := ufo.QueryMsgs(accountID, qparams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(pOut)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
		return
	}
}

// checkBoundsInt parses the int and clips it to the bounds, returns (val,ok)
func checkBoundsInt(in string, min, max int) (int, bool) {
	ret, err := strconv.Atoi(in)
	if err != nil || ret < min {
		return min, false
	} else if max < ret {
		return max, false
	}
	return ret, true
}

func parseIntField(in string, def, min, max int) int {
	ret, ok := checkBoundsInt(in, min, max)
	if !ok {
		return def
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
	folderIdx := chi.URLParam(r, "folderID")
	qp.FolderIdx = parseIntField(folderIdx, usecase.EnumInbox,
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
