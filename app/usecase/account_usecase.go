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

func (u *accountUsecase) GetSession() SessionUsecase {
	return u.session
}

func (u *accountUsecase) GetAccount(email string) (*Account, error) {
	acc, err := u.repo.Retrieve(email)
	if err != nil {
		return nil, ErrorNotFound
	}
	if acc == nil {
		return nil, ErrorNotFound
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
	return []*Account{}, ErrorNotFound
}

func (u *accountUsecase) RegisterAccount(email string) (*Account, error) {
	h := fnv.New64a()
	h.Write([]byte(email))
	uid := h.Sum64()
	if u.service.AlreadyExists(email) {
		return nil, fmt.Errorf("Account already exists")
	}
	acc := entity.NewAccount(entity.AccountIDType(uid), email)
	if err := u.repo.Create(acc); err != nil {
		return nil, err
	}
	out := toAccount([]*entity.Account{acc})
	return out[0], nil
}

func (u *accountUsecase) UpdateNameAccount(email string) error {
	h := fnv.New64a()
	h.Write([]byte(email))
	uid := h.Sum64()
	if u.service.AlreadyExists(email) {
		return fmt.Errorf("Account already exists")
	}
	Account := entity.NewAccount(entity.AccountIDType(uid), email)
	err := u.repo.Create(Account)
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
			ID:        strconv.FormatUint(uint64(account.GetID()), 16),
			Email:     account.GetEmail(),
			FirstName: account.GetFirstName(),
			LastName:  account.GetLastName(),
		}
	}
	return res
}
