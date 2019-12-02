package repo
import (
	"github.com/git-sim/tc/app/domain/entity"
)

type AccountRepo interface {
	Create(a *entity.Account) error
	Update(a *entity.Account) error
	Delete(a *entity.Account) error

	Retrieve(email string) (*entity.Account, error)
	RetrieveAll() ([]*entity.Account, error)
}