package repo

type StringRepo interface {
    Create(email, val string) error
    Update(email, val string) error
    Delete(email string) error

    Retrieve(email string) (string, error)
    RetrieveCount() (int, error)
    RetrieveAll() ([]string, error)
}