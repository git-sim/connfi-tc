package usecase

import (
	"fmt"
	"strconv"

	"github.com/git-sim/tc/app/domain/entity"
)

// AccountUsecase interface for account management
type AccountUsecase interface {
	GetAccountList() ([]*Account, error)
	GetAccount(email string) (*Account, error)
	RegisterAccount(email string) (*Account, error)
	UpdateNameAccount(email string, firstname *string, lastname *string) error
	DeleteAccount(email string) error

	GetSession() SessionUsecase
	IsRegisteredID(id string) bool
}

// An Account type for tranferring across the Usecase boundary
// provides isolation from details of entity.Account
type Account struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
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
