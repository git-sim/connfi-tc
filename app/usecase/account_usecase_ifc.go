package usecase

import (
	"errors"
)

// ErrorNotFound did not find account with the given email or ID
var ErrorNotFound = errors.New("Item not found")

// AccountUsecase interface for account management
type AccountUsecase interface {
	GetAccountList() ([]*Account, error)
	GetAccount(email string) (*Account, error)
	RegisterAccount(email string) (*Account, error)
	DeleteAccount(email string) error

	GetSession() SessionUsecase
}

// An Account type for tranferring across the Usecase boundary
// provides isolation from details of entity.Account
type Account struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}
