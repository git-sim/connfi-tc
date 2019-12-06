package ram

import (
    "fmt"
    "sync"
    "strconv"
    "sort"
    "github.com/git-sim/tc/app/domain/entity"
)

// Data type for Account in Ram, provides dependency inversion (isolation from) entity.Account 
type Account struct {
    ID    string
    Email string
}

// Impl of ram based account repository. Just a map[string]*Account
type accountRepo struct {
    mtx    *sync.Mutex
    accounts map[string]*Account
}

func NewAccountRepo() *accountRepo {
    return &accountRepo{
        mtx:    &sync.Mutex{},
        accounts: map[string]*Account{},
    }
}

func (r *accountRepo) Create(a *entity.Account) error {
    if(a == nil) {
        return fmt.Errorf("Invalid *entity.Account")
    }
    
    r.mtx.Lock()
    defer r.mtx.Unlock()

    r.accounts[a.GetEmail()] = &Account{
        ID:    string(a.GetID()),
        Email: a.GetEmail(),
    }
    return nil
}

func (r *accountRepo) Update(a *entity.Account) error {
    if(a == nil) {
        return fmt.Errorf("Invalid *entity.Account")
    }

    r.mtx.Lock()
    defer r.mtx.Unlock()

    if _ , ok := r.accounts[a.GetEmail()]; ok {
        r.accounts[a.GetEmail()] = &Account{
            ID:    string(a.GetID()),
            Email: a.GetEmail(),
        }
        return nil
    } else {
        return fmt.Errorf("Update error: entity.Account doesn't exist")
    }
}

func (r *accountRepo) Delete(a *entity.Account) error {
    if(a == nil) {
        return fmt.Errorf("Invalid entity.Account")
    }
    
    r.mtx.Lock()
    defer r.mtx.Unlock()
    delete(r.accounts, a.GetEmail())
    return nil
}

func (r *accountRepo) Retrieve(email string) (*entity.Account, error) {
    r.mtx.Lock()
    defer r.mtx.Unlock()

    account, ok := r.accounts[email]
    if ok {
        id, err := strconv.ParseInt(account.ID,10,64);
        if err == nil {
            return entity.NewAccount(entity.AccountID_t(id), account.Email), nil
        } else {
            return nil, err
        }
    }
    return nil, nil
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
    keys := make([]string,len(r.accounts))
    for k,_ := range r.accounts {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    var i int = 0
    for _,k := range keys {
	account := r.accounts[k]
        id, err := strconv.ParseInt(account.ID,10,64)
	if err == nil {
	    accounts[i] = entity.NewAccount(entity.AccountID_t(id), account.Email)
            i++
        }
    }
    return accounts, nil
}

