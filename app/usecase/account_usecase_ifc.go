package usecase

type AccountUsecase interface {
    GetAccountList() ([]*Account, error)
	GetAccount(email string) (*Account, error)
    RegisterAccount(email string) error
	DeleteAccount(email string) error
}

// An Account type for tranferring across the Usecase boundary
// provides isolation from details of entity.Account
type Account struct {
    ID    string
    Email string
}
