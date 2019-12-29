package usecase

// AccountUsecase interface for account management
type AccountUsecase interface {
	GetAccountList() ([]*Account, error)
	GetAccount(email string) (*Account, error)
	RegisterAccount(email string) (*Account, error)
	UpdateNameAccount(email string, firstname *string, lastname *string) error
	DeleteAccount(email string) error

	GetSession() SessionUsecase
}

// An Account type for tranferring across the Usecase boundary
// provides isolation from details of entity.Account
type Account struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}
