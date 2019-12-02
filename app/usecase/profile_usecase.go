package usecase
import (
    "github.com/git-sim/tc/app/domain/repo"
    "image"
    "fmt"
)


// Impl of name use case
const (
    firstNameEnum int = iota + 1
    lastNameEnum
)

type profileNameUsecase struct {
    whichName   int
    profileRepo repo.StringRepo
}

//Note internal function not exported
func newProfileNameUsecase(whichName int, nr repo.StringRepo) (*profileNameUsecase, error) {
    //Todo Is it bad etiquette for Newxxx() functions to return errors?
    if(nr==nil) {
        return nil, fmt.Errorf("Invalid repo.StringRepo")
    }
    u := &profileNameUsecase {
            whichName:      firstNameEnum,
            profileRepo:    nr,
    }
    return u,nil
}

func NewProfileFirstNameUsecase(nr repo.StringRepo) (*profileNameUsecase, error) {
    return newProfileNameUsecase(firstNameEnum, nr)
}

func NewProfileLastNameUsecase(nr repo.StringRepo) (*profileNameUsecase, error) {
    return newProfileNameUsecase(lastNameEnum, nr)
}
// ... NewProfileMiddleNameUsecase ...


func (u *profileNameUsecase) Set(email, val string) error {
    _ , err := u.profileRepo.Retrieve(email)
    if(err != nil) {
        // Create
        err = u.profileRepo.Create(email, val)
    } else {
        // Update
        err = u.profileRepo.Update(email, val)
    }
    return err
}

func (u *profileNameUsecase) Get(email string) (string, error) {
    currval , err := u.profileRepo.Retrieve(email)
    if(err != nil) {
        return "", err
    }
    return currval, nil
}

func (u *profileNameUsecase) GetCount() (int, error) {
    currval , err := u.profileRepo.RetrieveCount()
    if(err != nil) {
        return 0, err
    }
    return currval, nil
}

func (u *profileNameUsecase) GetList() ([]string, error) {
    currval , err := u.profileRepo.RetrieveAll()
    if(err != nil) {
        return []string{}, err
    }
    return currval, nil
}

// Impl of Profile Avatar Usecase   
type profileAvatarUsecase struct {
    profileRepo    repo.ImageRepo
}

func NewProfileAvatarUsecase(avr repo.ImageRepo) (*profileAvatarUsecase,error) {
    //Todo Is it bad etiquette for Newxxx() functions to return errors?
    if(avr==nil) {
        return nil, fmt.Errorf("Invalid repo.avatarRepo")
    }
    u := &profileAvatarUsecase {
        profileRepo: avr,
    }
    return u,nil
}

func (u *profileAvatarUsecase) Set(email string, val *image.Image) error {
    _ , err := u.profileRepo.Retrieve(email)
    if(err != nil) {
        // Create
        err = u.profileRepo.Create(email, val)
    } else {
        // Update
        err = u.profileRepo.Update(email, val)
    }
    return err
}

func (u *profileAvatarUsecase) Get(email string) (*image.Image, error) {
    currval , err := u.profileRepo.Retrieve(email)
    if(err != nil) {
        return nil, err
    }
    return currval, nil
}

func (u *profileAvatarUsecase) GetCount() (int, error) {
    currval , err := u.profileRepo.RetrieveCount()
    if(err != nil) {
        return 0, err
    }
    return currval, nil
}

func (u *profileAvatarUsecase) GetList() ([]*image.Image, error) {
    currval , err := u.profileRepo.RetrieveAll()
    if(err != nil) {
        return nil, err
    }
    return currval, nil
}

