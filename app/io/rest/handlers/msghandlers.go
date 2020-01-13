package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/git-sim/tc/app/usecase"
	"github.com/go-chi/chi"
)

// MessageCtxFunc returns a context handler to validate the message ID exists
func MessageCtxFunc(ufo usecase.FoldersUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlMessageID := chi.URLParam(r, "messageID")
			msgID, err := usecase.ToMsgID(urlMessageID)
			if err != nil {
				ReportUsecaseFault(w, err)
				return
			}

			v := r.Context().Value("dbAccountID")
			accountID, ok := v.(usecase.AccountIDType)
			if !ok {
				//context wasn't set right
				http.Error(w, "MessageCtxFunc AccountID context type assert failed",
					http.StatusInternalServerError)
				return
			}

			msg, err := ufo.GetOneMsg(accountID, msgID)
			if err != nil && msg != nil {
				fmt.Printf("AccountCtxFunc Message Not found AccountID %d MessageID Not found %s\n", accountID, urlMessageID)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			// Save away the message and the message ID for this request so
			ctx := context.WithValue(r.Context(), "dbMsgID", msgID)
			ctx2 := context.WithValue(ctx, "message", msg)
			next.ServeHTTP(w, r.WithContext(ctx2))
		})
	}
}

// CreateMessage ...
func CreateMessage(mu usecase.MsgUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode
		var inmsg usecase.IngressMsg
		err := decodeJSONBody(w, r, &inmsg)
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

		//Enq the message
		outmsgid, err := mu.EnqueueMsg(&inmsg)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		// Returning the ingested message in the response, could just
		//    return the id for brevity. The client can get the whole message if they want it.
		outmsg, err := mu.RetrieveMsg(outmsgid)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(&outmsg)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// GetMessages ...
//func GetMessages(mu usecase.MsgUsecase, ufo usecase.FoldersUsecase, u usecase.AccountUsecase) http.HandlerFunc {
//
//}

// GetMessage ... this could just be a HandlerFunc but keeping it for symmetry and if we need to add params
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The message is already cached in the context for this request
		// extract and encode it
		outmsg, ok := getMsgFromContext(w, r)
		if !ok {
			return
		}
		err := json.NewEncoder(w).Encode(&outmsg)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// PutMessage ...
func PutMessage(ufo usecase.FoldersUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The message is already cached in the context for this request
		// extract, modify, update  (concurrency?)
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}

		msgID, ok := getMsgIDFromContext(w, r)
		if !ok {
			return
		}

		// Decode body
		var inmsg usecase.MsgEntry
		err := decodeJSONBody(w, r, &inmsg)
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

		err = ufo.UpdateMsg(accountID, msgID, inmsg)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		outmsg, err := ufo.GetOneMsg(accountID, msgID)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(outmsg)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
}

// DeleteMessage ...
func DeleteMessage(ufo usecase.FoldersUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, ok := getAccountIDFromContext(w, r)
		if !ok {
			return
		}

		msgID, ok := getMsgIDFromContext(w, r)
		if !ok {
			return
		}

		outmsg, ok := getMsgFromContext(w, r)
		if !ok {
			return
		}

		// finally we're ready to delete
		err := ufo.DeleteMsg(accountID, msgID)
		if err != nil {
			ReportUsecaseFault(w, err)
			return
		}

		// Sendback the message that's been deleted
		err = json.NewEncoder(w).Encode(&outmsg)
		if err != nil {
			ReportJSONFault(w, err)
			return
		}
	}
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
