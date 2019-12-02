package entity

// Account represents an messaging account. Email should be unique
type AccountID_t uint64
type Account struct {
    id        AccountID_t
    email     string    
    FirstName string
    LastName  string
}

//How to create a new account verifying the email is unique
// We need a collection of existing accounts.
// When creating a new account need to check that the email is unique
// When created the account will get assigned an inbox (collection of inbox messages)
// Inbox could start empty or with initial welcome message (from admin).


// NewAccount instantiates a new Account. 
func NewAccount(newid AccountID_t, newemail string) *Account {
    return &Account {
        id:        newid,
        email:     newemail,
    }
}

func (a *Account) GetID() AccountID_t {
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

// NewAccounts instanties a slice of new Accounts
func NewAccounts(Emails ...string) []*Account {
    var as []*Account
    for i, email := range Emails {
        as = append(as, NewAccount(AccountID_t(i+1), email))
    }
    return as
}