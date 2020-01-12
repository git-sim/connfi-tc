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
func (u *accountUsecase) RegisterAccount(account *Account) (*Account, error) {
	if account == nil {
		return nil, NewEs(EsArgInvalid, "User Account")
	}

	if u.service.AlreadyExists(account.Email) {
		return nil, NewEs(EsAlreadyExists, "User Account email")
	}

	// Create the account and associated structures in the system
	//   A Delete account should undo the below in reverse order to make sure
	//   we have a good cleanup
	uid := GetUID(account.Email)
	acc := entity.NewAccount(entity.AccountIDType(uid), account.Email)
	acc.FirstName = account.FirstName
	acc.LastName = account.LastName
	if err := u.repo.Create(acc); err != nil {
		return nil, err
	}
	u.service.NotifyRegisterAccount(*acc)
	out := toAccount([]*entity.Account{acc})
	return out[0], nil
}

// RegisterAccountByEmail this is one of the major events in the system creating the structures needed for the account.
func (u *accountUsecase) RegisterAccountByEmail(email string) (*Account, error) {
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

func (u *accountUsecase) GetAccount(id string) (*Account, error) {
	accID, err := ToAccountID(id)
	if err != nil {
		return nil, err
	}

	acc, err := u.repo.RetrieveByID(entity.AccountIDType(accID))
	if err != nil {
		return nil, NewEs(EsNotFound, "id")
	}
	if acc == nil {
		return nil, NewEs(EsEmpty, "User Account")
	}
	out := toAccount([]*entity.Account{acc})
	return out[0], nil
}

func (u *accountUsecase) UpdateAccount(account *Account) error {
	// Check that it exists, if so validate and update
	accountID, err := ToAccountID(account.ID)
	if err != nil {
		return err
	}

	// Initialize the next value with the current value
	nextAccount, err := u.repo.RetrieveByID(entity.AccountIDType(accountID))
	if err != nil {
		return err
	}

	nextAccount.FirstName = account.FirstName
	nextAccount.LastName = account.LastName

	err = u.repo.Update(nextAccount)
	return err
}

func (u *accountUsecase) DeleteAccount(id string) error {
	accID, err := ToAccountID(id)
	if err != nil {
		// Problem converting from string to id, details are in err
		return err
	}

	ok := u.IsRegisteredID(id)
	if !ok {
		// meets the contract that after this call the account doesn't exist
		return nil
	}

	// todo don't allow deleting admin no matter what
	acc, err := u.repo.RetrieveByID(entity.AccountIDType(accID))
	u.service.NotifyDeleteAccount(*acc)
	err = u.repo.Delete(acc)
	return err
}

func (u *accountUsecase) GetAccountByEmail(email string) (*Account, error) {
	acc, err := u.repo.RetrieveByEmail(email)
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
	a, err := u.repo.RetrieveByEmail(email)
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

func (u *accountUsecase) DeleteAccountByEmail(email string) error {
	a, err := u.repo.RetrieveByEmail(email)
	if err != nil {
		return err
	}
	if a != nil {
		u.repo.Delete(a)
	}
	return nil
}

func (u *accountUsecase) IsRegisteredID(id string) bool {
	accID, err := ToAccountID(id)
	if err != nil {
		return false
	}
	return u.service.AlreadyExistsByID(entity.AccountIDType(accID))
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
