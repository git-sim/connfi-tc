package repo

type GenericKeyT uint64
type Generic interface {
	Create(id GenericKeyT, val interface{}) error
	Update(id GenericKeyT, val interface{}) error
	Delete(id GenericKeyT) error

	Retrieve(id GenericKeyT) (interface{}, error)
	RetrieveCount() (int, error)
	RetrieveAll() ([]interface{}, error) //todo maybe should return map[uint64]*string so we have id and value info
}
