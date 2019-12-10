package repo

type StringRepo interface {
    Create(id uint64, val string) error
    Update(id uint64, val string) error
    Delete(id uint64) error

    Retrieve(id uint64) (string, error)
    RetrieveCount() (int, error)
    RetrieveAll() ([]*string, error) //todo maybe should return map[uint64]*string so we have id,string info
}