package entity

// Account represents an messaging account. Email should be unique
type Account struct {
	id        AccountIDType
	email     string
	FirstName string
	LastName  string
}

// AccountIDType specifies the id type
type AccountIDType uint64

const AccountIDBits = 64
const AccountIDStringBase = 16

//How to create a new account verifying the email is unique
// We need a collection of existing accounts.
// When creating a new account need to check that the email is unique
// When created the account will get assigned an inbox (collection of inbox messages)
// Inbox could start empty or with initial welcome message (from admin).

// NewAccount instantiates a new Account.
func NewAccount(newid AccountIDType, newemail string) *Account {
	return &Account{
		id:    newid,
		email: newemail,
	}
}

func (a *Account) GetID() AccountIDType {
	return a.id
}

func (a *Account) GetEmail() string {
	return a.email
}

func (a *Account) GetFirstName() string {
	return a.FirstName
}

func (a *Account) GetLastName() string {
	return a.LastName
}

// NewAccounts instantiates a slice of new Accounts
func NewAccounts(Emails ...string) []*Account {
	var as []*Account
	for i, email := range Emails {
		as = append(as, NewAccount(AccountIDType(i+1), email))
	}
	return as
}
