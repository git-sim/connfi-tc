package repo

import (
	"image"
)

// Basic CRUD ifc to a PublicProfile

type ProfileRepo interface {
	Create(id uint64, val *PublicProfile) error
	Update(id uint64, val *PublicProfile) error
	Delete(id uint64) error

	Retrieve(id uint64) (*PublicProfile, error)
	RetrieveCount() (int, error)
}

const (
	EnumFirstName = iota
	EnumLastName
	EnumSalutation //Mr., Ms., Dr., Lt Col., etc
	EnumSuffix     //Jr., Sr., III, PhD, Esq, etc
	EnumBio
	EnumNumProfileStringFields
)
const (
	EnumAvatar = iota
	EnumBackground
	EnumNumProfileImageFields
)

type PublicProfile struct {
	NameAndBios [EnumNumProfileStringFields]string
	Pics        [EnumNumProfileImageFields]*image.Image
	//... others
}
