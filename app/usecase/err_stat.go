package usecase

import "fmt"

//This module defines generally used errors and statuses in the usecase packages
// Have the basic ErrStat type and some commonly error codes in the usecases/interactors.
// Each usecase can extend it by adding their own statuses in a separate file
// Mainly need a simple way to detect a specific error code and be able to transform it (ie to httpStatus) or take action,
// instead of just passing it up the chain.

// ErrStat error and status type that meets the Error interface
type ErrStat struct {
	Code int
	Msg  string
}

// Predefined codes the numbers themselves aren't significant need to be able to compare them
const (
	// MAINT NOTE if adding removing errors remember to update the string map below
	EsOk = 0

	// Pos numbers for regular operational errors
	EsEmpty           = 101
	EsAlreadyExists   = 102
	EsArgInvalid      = 201
	EsArgConvFail     = 202
	EsForbidden       = 301
	EsNotFound        = 302
	EsAlreadyReported = 401

	// Convention use negative numbers for faults and internal issues
	EsOutOfResources = -100
	EsInternalError  = -500
	EsNotImplemented = -501

	EsMaxBaseError = 1000 //to allow other modules to start at a higher value
	EsMinBaseError = -1000
)

var esText = map[int]string{
	EsOk: "Ok",

	// Pos numbers for regular operational errors
	EsEmpty:           "Empty",
	EsAlreadyExists:   "Already Exists",
	EsArgInvalid:      "Arg Invalid",
	EsArgConvFail:     "Arg Conversion Fail",
	EsForbidden:       "Fobidden ",
	EsNotFound:        "Not Found",
	EsAlreadyReported: "Already Reported",

	// Convention use negative numbers for faults and internal issues
	EsOutOfResources: "Out of Resources",
	EsInternalError:  "Internal Error",
	EsNotImplemented: "Not Implemented",

	EsMaxBaseError: "MaxBase",
	EsMinBaseError: "MinBase",
}

// ErrStatusText returns a text for the Status codes. It returns the empty
// string if the code is unknown.
func ErrStatusText(code int) string {
	return esText[code]
}

// Error function to make ErrStat meet the Error ifc
func (es *ErrStat) Error() string {
	return fmt.Sprintf("ErrStat %d: %s %s",
		es.Code, ErrStatusText(es.Code), es.Msg)
}

// NewEs Helper functions same as &ErrStat{EsCode,"explain"}
func NewEs(code int, msg string) *ErrStat {
	// could enable logging here for debug
	return &ErrStat{code, msg}
}
