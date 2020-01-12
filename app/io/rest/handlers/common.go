package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/git-sim/tc/app/usecase"
)

// setupCORS Cross Origin request
func setupCORS(w http.ResponseWriter, r *http.Request) {
	//
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// SetupCORSHandler for use as a middleware
func SetupCORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		next.ServeHTTP(w, r)
	})
}

//The following method is from alexedwards.net/blog/how-to-properly-parse-a-json-request-body
type ErrorJSONDecode struct {
	status int
	msg    string
}

func (mr *ErrorJSONDecode) Error() string {
	return mr.msg
}

// MaxJsonBodySize set a max limit, for how much data we'll accept in a request
const MaxJSONBodySize = 1 * 1024 * 1024

// decodeJSONBody for decoding body params (in POST's and PUT's)
func decodeJSONBody(w http.ResponseWriter, r *http.Request, out interface{}) error {

	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize) //limit how large a json body we'll handle
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields
	err := d.Decode(&out)

	// Reference: The error handling from a json Decode is heavily leveraged from
	// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", MaxJSONBodySize)
			return &ErrorJSONDecode{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	if d.More() {
		msg := "Request body must only contain a single JSON object"
		return &ErrorJSONDecode{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

// getAccIDFromSession checks the session cookie returns the session info(accIDString, ok, aut)
func getAccIDFromSession(u usecase.AccountUsecase, r *http.Request) (accIDString string, ok bool, auth bool) {

	accIDString = ""
	ok = false
	auth = false

	//Always returns a session
	session, _ := u.GetSession().FromReq(r)
	// Could do auth here, we're interested in getting the AccountId of the user
	if auth, ok = session.Values["authenticated"].(bool); ok && auth {
		accIDString, ok = session.Values["id"].(string)
	}

	// Try from url params, should be able to get rid of the above
	if ok == false {
		r.ParseForm()
		accIDString = r.FormValue("accid")
		if u.IsRegisteredID(accIDString) {
			ok = true
			auth = true
		}
	}
	return accIDString, ok, auth
}

// ErrStatToHTTPCode translate from internal code to best http code. If no match returns http.StatusInternalServerError
func ErrStatToHTTPCode(es *usecase.ErrStat) int {
	out := http.StatusInternalServerError
	if es == nil {
		return out
	}
	in := es.Code

	switch in {
	case usecase.EsOk:
		out = http.StatusOK
	case usecase.EsEmpty:
		out = http.StatusNoContent
	case usecase.EsAlreadyExists:
		out = http.StatusFound
	case usecase.EsArgInvalid:
		out = http.StatusBadRequest
	case usecase.EsArgConvFail:
		out = http.StatusBadRequest
	case usecase.EsForbidden:
		out = http.StatusForbidden
	case usecase.EsNotFound:
		out = http.StatusNotFound
	case usecase.EsAlreadyReported:
		out = http.StatusAlreadyReported
	case usecase.EsOutOfResources:
		out = http.StatusInternalServerError
	case usecase.EsInternalError:
		out = http.StatusInternalServerError
	case usecase.EsNotImplemented:
		out = http.StatusInternalServerError
	default:
		out = http.StatusInternalServerError
	}
	return out
}

// ReportJSONFault helper function to report an error with JSON, and do any debug logging to alert devs
func ReportJSONFault(w http.ResponseWriter, err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("Debug file: %s line: %d Error: %s", file, line, err.Error())

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
