package usecase

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/git-sim/tc/app/domain/repo"
)

type profileUsecase struct {
	profileRepo repo.ProfileRepo
}

func NewProfileUsecase(pr repo.ProfileRepo) ProfileUsecase {
	// should be an assert
	if pr == nil {
		log.Fatal(fmt.Errorf("invalid repo.ProfileRepo"))
	}
	u := &profileUsecase{profileRepo: pr}
	return u
}

func (pu *profileUsecase) Set(id uint64, val *repo.PublicProfile) error {
	// Update
	return pu.profileRepo.Update(id, val)
}

func (pu *profileUsecase) Get(id uint64) (*repo.PublicProfile, error) {
	return pu.profileRepo.Retrieve(id)
}

func (pu *profileUsecase) GetCount() (int, error) {
	return pu.profileRepo.RetrieveCount()
}

func (pu *profileUsecase) Delete(id uint64) error {
	return pu.profileRepo.Delete(id)
}

func (pu *profileUsecase) CreateDefaultProfile(id uint64) error {
	p := repo.PublicProfile{}
	p.NameAndBios[repo.EnumBio] = "No Bio Recorded"
	defaultAvatar, err := os.Open("./testdata/img1.png")
	if err == nil {
		defer defaultAvatar.Close()
		avImage, err := png.Decode(defaultAvatar)
		if err == nil {
			p.Pics[repo.EnumAvatar] = &avImage
		}
	}

	defaultBg, err := os.Open("./testdata/img2.png")
	if err == nil {
		defer defaultBg.Close()
		bgImage, err := png.Decode(defaultBg)
		if err == nil {
			p.Pics[repo.EnumBackground] = &bgImage
		}
	}
	return pu.profileRepo.Create(id, &p)
}
