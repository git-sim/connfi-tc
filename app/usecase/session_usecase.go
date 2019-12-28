package usecase

import (
	"net/http"

	"github.com/git-sim/tc/app/domain/service"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	//TODO In a real project this needs be an environment variable, can retrieve via os.Getenv("SESSION_KEY")
	defaultkey = []byte("super-secret-key")
)

// Impl of LoginUsecase interface
type sessionUsecase struct {
	cookieStore *sessions.CookieStore
	service     *service.AccountService
}

// NewSessionUsecase - the sessionkey is the seed for the session cookie store, service is he AccountService
func NewSessionUsecase(sessionkey []byte, service *service.AccountService) SessionUsecase {

	key := defaultkey
	if sessionkey != nil {
		key = sessionkey
	}

	return &sessionUsecase{
		cookieStore: sessions.NewCookieStore(key),
		service:     service,
	}
}

// SessionUsecase.FreomReq returns the session cookie
func (ul *sessionUsecase) FromReq(r *http.Request) (*sessions.Session, error) {
	return ul.cookieStore.Get(r, "session-cookie")
}
