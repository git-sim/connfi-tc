package usecase

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
)

// accountUsecase impl
type accountUsecase struct {
	repo    repo.AccountRepo
	session SessionUsecase
	service *service.AccountService
}

// NewAccountUsecase - repo is the interface for the Account Repository (db Or in memory)
func NewAccountUsecase(repo repo.AccountRepo, session SessionUsecase, service *service.AccountService) AccountUsecase {
	return &accountUsecase{
		repo:    repo,
		session: session,
		service: service,
	}
}

// RegisterAccount this is one of the major events in the system creating the structures needed for the account.
func (u *accountUsecase) RegisterAccount(email string) (*Account, error) {
	if u.service.AlreadyExists(email) {
		return nil, NewEs(EsAlreadyExists, "User Account")
	}

	// Create the account and associated structures in the system
	//   A Delete account should undo the below in reverse order to make sure
	//   we have a good cleanup
	uid := GetUID(email)
	acc := entity.NewAccount(entity.AccountIDType(uid), email)
	if err := u.repo.Create(acc); err != nil {
		return nil, err
	}
	u.service.NotifyRegisterAccount(*acc)
	out := toAccount([]*entity.Account{acc})
	return out[0], nil
}

func (u *accountUsecase) GetSession() SessionUsecase {
	return u.session
}

func (u *accountUsecase) GetAccount(email string) (*Account, error) {
	acc, err := u.repo.Retrieve(email)
	if err != nil {
		return nil, NewEs(EsNotFound, "Email")
	}
	if acc == nil {
		return nil, NewEs(EsEmpty, "User Account")
	}
	out := toAccount([]*entity.Account{acc})
	return out[0], nil
}

func (u *accountUsecase) GetAccountList() ([]*Account, error) {
	Accounts, err := u.repo.RetrieveAll()
	if err != nil {
		return nil, err
	}
	if len(Accounts) > 0 {
		return toAccount(Accounts), nil
	}
	return []*Account{}, NewEs(EsEmpty, "Accounts")
}

func (u *accountUsecase) UpdateNameAccount(email string, firstname *string, lastname *string) error {
	a, err := u.repo.Retrieve(email)
	if err != nil {
		return err
	}
	if a == nil {
		return NewEs(EsEmpty,
			fmt.Sprintf("Account with email %s", email))
	}

	// leave email and id the same just update the names if they exist
	if firstname != nil {
		a.FirstName = *firstname
	}
	if lastname != nil {
		a.LastName = *lastname
	}
	err = u.repo.Update(a)
	return err
}

func (u *accountUsecase) DeleteAccount(email string) error {
	a, err := u.repo.Retrieve(email)
	if err != nil {
		return err
	}
	if a != nil {
		u.repo.Delete(a)
	}
	return nil
}

// Conversion function from entity.Account to usecase.Account
func toAccount(Accounts []*entity.Account) []*Account {
	res := make([]*Account, len(Accounts))
	for i, account := range Accounts {
		res[i] = &Account{
			ID:        strconv.FormatUint(uint64(account.GetID()), entity.AccountIDStringBase),
			Email:     account.GetEmail(),
			FirstName: account.GetFirstName(),
			LastName:  account.GetLastName(),
		}
	}
	return res
}

// details of getting a unique id for the account
func GetUID(in string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(in))
	uid := h.Sum64()
	return uid
}
