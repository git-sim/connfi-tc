package usecase
import (
    "errors"
    "fmt"
    "github.com/git-sim/tc/app/domain/entity"
    "github.com/git-sim/tc/app/domain/repo"
    "github.com/git-sim/tc/app/domain/service"
    "hash/fnv"
    "strconv"
)

// Impl of AccountUseCase interface
type accountUsecase struct {
    repo    repo.AccountRepo
    service *service.AccountService
}
var errNotFound = errors.New("Item not found")

func NewAccountUsecase(repo repo.AccountRepo, service *service.AccountService) *accountUsecase {
    return &accountUsecase{
        repo:    repo,
        service: service,
    }
}

func (u *accountUsecase) GetAccount(email string) (*Account, error) {
    acc, err := u.repo.Retrieve(email)
    if err != nil {
        return nil, err
    }
    if acc == nil {
        return nil, errNotFound
    }
    out := toAccount([]*entity.Account{acc})
    return out[0], nil
}

func (u *accountUsecase) GetAccountList() ([]*Account, error) {
    Accounts, err := u.repo.RetrieveAll()
    if err != nil {
        return nil, err
    }
    if(len(Accounts)>0) {
        return toAccount(Accounts), nil
    }
    return []*Account{}, errNotFound
}

func (u *accountUsecase) RegisterAccount(email string) error {
    h := fnv.New64a()
    h.Write([]byte(email))
    uid := h.Sum64()
    if err := u.service.AlreadyExists(email); err == nil {
        return fmt.Errorf("Account already exists")
    }
    Account := entity.NewAccount(entity.AccountID_t(uid), email)
    if err := u.repo.Create(Account); err != nil {
        return err
    }
    return nil
}

func (u *accountUsecase) UpdateNameAccount(email string) error {
    h := fnv.New64a()
    h.Write([]byte(email))
    uid := h.Sum64()
    if err := u.service.AlreadyExists(email); err == nil {
        return fmt.Errorf("Account already exists")
    }
    Account := entity.NewAccount(entity.AccountID_t(uid), email)
    if err := u.repo.Create(Account); err != nil {
        return err
    }
    return nil
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
            ID:    strconv.FormatUint(uint64(account.GetID()),16),
            Email: account.GetEmail(),
            FirstName: account.GetFirstName(),
            LastName: account.GetLastName(),
        }
    }
    return res
}

