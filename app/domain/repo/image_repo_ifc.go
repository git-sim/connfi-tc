package repo
import ( 
	"image"
)

type ImageRepo interface {
	Create(email string, val *image.Image) error
	Update(email string, val *image.Image) error
	Delete(email string) error

	Retrieve(email string) (*image.Image, error)
	RetrieveCount() (int, error)
	RetrieveAll() ([]*image.Image, error)
}