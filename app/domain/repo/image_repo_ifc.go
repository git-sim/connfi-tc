package repo
import ( 
    "image"
)

type ImageRepo interface {
    Create(id uint64, val *image.Image) error
    Update(id uint64, val *image.Image) error
    Delete(id uint64) error

    Retrieve(id uint64) (*image.Image, error)
    RetrieveCount() (int, error)
    RetrieveAll() ([]*image.Image, error) //todo maybe should return map[uint64]*string so we have id,string info
}