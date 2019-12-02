package usecase
import (
	"image"
)

// One interface for first and last name
type ProfileNameUsecase interface {
	Set(email, val string) error
	Get(email string) (string, error)
	GetCount() (int, error)
	GetList()  ([]string, error)
}

// Very similar to above just differs on target value type	
type ProfileAvatarUsecase interface {
	Set(email string, val *image.Image) error
	Get(email string) (*image.Image, error)
	GetCount() (int, error)
	GetList()  ([]*image.Image, error)
}

