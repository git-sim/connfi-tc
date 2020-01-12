package usecase

import (
	"fmt"
	"strconv"

	"github.com/git-sim/tc/app/domain/entity"
)

// AccountUsecase interface for account management
type AccountUsecase interface {
	RegisterAccount(account *Account) (*Account, error)
	RegisterAccountByEmail(email string) (*Account, error)
	GetAccountList() ([]*Account, error)
	GetAccount(id string) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id string) error

	// deprecated
	UpdateNameAccount(email string, firstname *string, lastname *string) error
	DeleteAccountByEmail(email string) error
	// /deprecated

	GetSession() SessionUsecase
	IsRegisteredID(id string) bool
}

// An Account type for tranferring across the Usecase boundary
// provides isolation from details of entity.Account
type Account struct {
	ID        string `json:"accid,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type AccountIDType entity.AccountIDType

const AccountIDBits = entity.AccountIDBits
const AccounIDStringBase = entity.AccountIDStringBase

// Conversion functions

// AccountIDToString ...
func AccountIDToString(id AccountIDType) string {
	return strconv.FormatUint(uint64(id), AccounIDStringBase)
}

// ToAccountID so we're all on the same format
func ToAccountID(idString string) (AccountIDType, error) {
	id, err := strconv.ParseUint(idString, AccounIDStringBase, AccountIDBits)
	if err != nil {
		return AccountIDType(0), NewEs(EsArgConvFail,
			fmt.Sprintf("idString %s", idString))
	}
	return AccountIDType(id), nil
}
