package ram

import (
	"sort"
	"strconv"
	"sync"

	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/usecase"
)

// Data type for Account in Ram, provides dependency inversion (isolation from) entity.Account
type Account struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

// Account.toEntityAccount conversion helper
func (ra *Account) toEntityAccount(id64 entity.AccountIDType) *entity.Account {
	ret := entity.NewAccount(id64, ra.Email)
	ret.FirstName = ra.FirstName
	ret.LastName = ra.LastName
	return ret
}

// Impl of ram based account repository. Just a map[string]*Account
type accountRepo struct {
	mtx      *sync.Mutex
	accounts map[entity.AccountIDType]*Account
}

func NewAccountRepo() *accountRepo {
	return &accountRepo{
		mtx:      &sync.Mutex{},
		accounts: map[entity.AccountIDType]*Account{},
	}
}

func (r *accountRepo) Create(a *entity.Account) error {
	if a == nil {
		return usecase.NewEs(usecase.EsArgInvalid, "*entity.Account")
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.accounts[a.GetID()] = &Account{
		ID:        GetIDString(a.GetID()),
		Email:     a.GetEmail(),
		FirstName: a.GetFirstName(),
		LastName:  a.GetLastName(),
	}
	return nil
}

func (r *accountRepo) Update(a *entity.Account) error {
	if a == nil {
		return usecase.NewEs(usecase.EsArgInvalid, "*entity.Account")
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.accounts[a.GetID()]; ok {
		r.accounts[a.GetID()] = &Account{
			ID:        GetIDString(a.GetID()),
			Email:     a.GetEmail(),
			FirstName: a.GetFirstName(),
			LastName:  a.GetLastName(),
		}
		return nil
	} else {
		return usecase.NewEs(usecase.EsNotFound, "entity.Account")
	}
}

func (r *accountRepo) Delete(a *entity.Account) error {
	if a == nil {
		return usecase.NewEs(usecase.EsArgInvalid, "*entity.Account")
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.accounts, a.GetID())
	return nil
}

func (r *accountRepo) RetrieveByEmail(email string) (*entity.Account, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for k, v := range r.accounts {
		if v.Email == email {
			ret := v.toEntityAccount(k)
			return ret, nil
		}
	}
	return nil, usecase.NewEs(usecase.EsNotFound, "Email")

}

func (r *accountRepo) RetrieveByID(id entity.AccountIDType) (*entity.Account, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	v, ok := r.accounts[id]
	if ok {
		ret := v.toEntityAccount(id)
		return ret, nil
	}
	return nil, usecase.NewEs(usecase.EsNotFound, "account id")
}

func (r *accountRepo) RetrieveCount() (int, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return len(r.accounts), nil
}

func (r *accountRepo) RetrieveAll() ([]*entity.Account, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	accounts := make([]*entity.Account, len(r.accounts))
	type pairIdEmail struct {
		AcctID entity.AccountIDType
		Email  string
	}
	idEmails := make([]pairIdEmail, len(r.accounts))
	var i int = 0
	for k, v := range r.accounts {
		idEmails[i] = pairIdEmail{AcctID: k, Email: v.Email}
		i++
	}
	sort.Slice(idEmails, func(i, j int) bool {
		return idEmails[i].Email < idEmails[j].Email
	})

	var j int = 0
	for _, idEmail := range idEmails {
		account := r.accounts[idEmail.AcctID]
		accounts[j] = account.toEntityAccount(idEmail.AcctID)
		j++
	}
	return accounts, nil
}

// GetIDString Helpers for conversions
func GetIDString(id entity.AccountIDType) string {
	return strconv.FormatUint(uint64(id), entity.AccountIDStringBase)
}

// The below is provided for symmetry,reference but not actually used so far.
// How to convert from str to id with the correct errors
//func fromStrToID(s string) (entity.AccountIDType, error) {
//	n, err := strconv.ParseUint(s, entity.AccountIDStringBase, entity.AccountIDBits)
//	if err != nil {
//		return 0, usecase.NewEs(usecase.EsArgInvalid,
//			fmt.Sprintf("bad ID string %s, must be uint64 encoded in hex", s))
//	}
//	return entity.AccountIDType(n), nil
//}
