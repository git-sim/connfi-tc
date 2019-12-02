package usecase
import (
	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
	"hash/fnv"
)

// Impl of AccountUseCase interface
type accountUsecase struct {
	repo    repo.AccountRepo
	service *service.AccountService
}

func NewAccountUsecase(repo repo.AccountRepo, service *service.AccountService) *accountUsecase {
	return &accountUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *accountUsecase) GetAccountList() ([]*Account, error) {
	Accounts, err := u.repo.RetrieveAll()
	if err != nil {
		return nil, err
	}
	return toAccount(Accounts), nil
}

func (u *accountUsecase) RegisterAccount(email string) error {
	h := fnv.New64a()
	h.Write([]byte(email))
	uid := h.Sum64()
	if err := u.service.AlreadyExists(email); err != nil {
		return err
	}
	Account := entity.NewAccount(entity.AccountID_t(uid), email)
	if err := u.repo.Create(Account); err != nil {
		return err
	}
	return nil
}

// Conversion function from entity.Account to usecase.Account
func toAccount(Accounts []*entity.Account) []*Account {
	res := make([]*Account, len(Accounts))
	for i, account := range Accounts {
		res[i] = &Account{
			ID:    string(account.GetID()),
			Email: account.GetEmail(),
		}
	}
	return res
}

