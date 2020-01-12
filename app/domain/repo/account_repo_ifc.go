package repo

import (
	"github.com/git-sim/tc/app/domain/entity"
)

type AccountRepo interface {
	Create(a *entity.Account) error
	Update(a *entity.Account) error
	Delete(a *entity.Account) error

	RetrieveByEmail(email string) (*entity.Account, error)
	RetrieveByID(id entity.AccountIDType) (*entity.Account, error)
	RetrieveCount() (int, error)
	RetrieveAll() ([]*entity.Account, error)
}
