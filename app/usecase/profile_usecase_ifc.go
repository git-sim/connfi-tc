package usecase
import (
    "image"
)

// One interface for first and last name
type ProfileStringUsecase interface {
    Set(id uint64, val string) error
    Get(id uint64) (string, error)
    GetCount() (int, error)
    GetList()  ([]*string, error)
}

// Very similar to above just differs on target value type  
type ProfileImageUsecase interface {
    Set(id uint64, val *image.Image) error
    Get(id uint64) (*image.Image, error)
    GetCount() (int, error)
    GetList()  ([]*image.Image, error)
}

