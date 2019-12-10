package usecase
import (
    "fmt"
    "github.com/git-sim/tc/app/domain/repo"
    "image"
    "log"
)

type profileStringUsecase struct {
    profileRepo repo.StringRepo
}

//Note internal function not exported
func NewProfileStringUsecase(nr repo.StringRepo) *profileStringUsecase {
    // should be an assert
    if nr==nil {
        log.Fatal(fmt.Errorf("invalid repo.StringRepo")) 
    }
    u := &profileStringUsecase {profileRepo: nr }
    return u
}

func (u *profileStringUsecase) Set(id uint64, val string) error {
    _ , err := u.profileRepo.Retrieve(id)
    if err != nil {
        // Create
        err = u.profileRepo.Create(id, val)
    } else {
        // Update
        err = u.profileRepo.Update(id, val)
    }
    return err
}

func (u *profileStringUsecase) Get(id uint64) (string, error) {
    val , err := u.profileRepo.Retrieve(id)
    if err != nil {
        return "", err
    }
    return val, nil
}

func (u *profileStringUsecase) GetCount() (int, error) {
    count , err := u.profileRepo.RetrieveCount()
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (u *profileStringUsecase) GetList() ([]*string, error) {
    currval , err := u.profileRepo.RetrieveAll()
    if err != nil {
        return []*string{}, err
    }
    return currval, nil
}

// Impl of Profile Image Usecase
type profileImageUsecase struct {
    profileRepo    repo.ImageRepo
}

func NewProfileImageUsecase(ir repo.ImageRepo) *profileImageUsecase {
    // should be a compile time assert
    if ir==nil {
        log.Fatal(fmt.Errorf("invalid repo.ImageRepo"))
    }
    u := &profileImageUsecase {
        profileRepo: ir,
    }
    return u
}

func (u *profileImageUsecase) Set(id uint64, val *image.Image) error {
    _ , err := u.profileRepo.Retrieve(id)
    if err != nil  {
        // Create
        err = u.profileRepo.Create(id, val)
    } else {
        // Update
        err = u.profileRepo.Update(id, val)
    }
    return err
}

func (u *profileImageUsecase) Get(id uint64) (*image.Image, error) {
    currval , err := u.profileRepo.Retrieve(id)
    if err != nil {
        return nil, err
    }
    return currval, nil
}

func (u *profileImageUsecase) GetCount() (int, error) {
    currval , err := u.profileRepo.RetrieveCount()
    if err != nil {
        return 0, err
    }
    return currval, nil
}

func (u *profileImageUsecase) GetList() ([]*image.Image, error) {
    currval , err := u.profileRepo.RetrieveAll()
    if err != nil {
        return nil, err
    }
    return currval, nil
}

