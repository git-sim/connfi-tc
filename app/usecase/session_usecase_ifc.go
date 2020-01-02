package usecase

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionUsecase Interface -
type SessionUsecase interface {
	
	FromReq(r *http.Request) (*sessions.Session, error)
}
