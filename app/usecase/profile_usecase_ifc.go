package usecase

import (
	"github.com/git-sim/tc/app/domain/repo"
)

// The profile elements don't need to be fancy just a crud ifc
type ProfileUsecase interface {
	Set(id uint64, val *repo.PublicProfile) error
	Get(id uint64) (*repo.PublicProfile, error)
	GetCount() (int, error)
	Delete(id uint64) error

	CreateDefaultProfile(id uint64) error
}
