package ram

import (
    "fmt"
    "image"
    _ "image/png"
    "log"
    "sync"
)

// Impl of ram based profile repository
type profileRepo struct {
    mtx    *sync.Mutex
    Profiles   map[uint64]PublicProfile
}

func NewProfileRepo() *profileRepo {
    return &profileRepo {
        mtx: &sync.Mutex{},
        Profiles: map[uint64]PublicProfile{},
    }
}

// Gets notified that a field from the public profile has been removed
// clears out the entry in the map, if all fields are removed
func (pr *profileRepo) DeleteNotify(id uint64) {
    pubProfile, ok := pr.Profiles[id]
    if ok {
        EmptyProfile := PublicProfile{}
        if pubProfile == EmptyProfile {
            delete(pr.Profiles, id)
        }
    }
}

const (
    EnumFirstName = iota
    EnumLastName
    EnumSalutation  //Mr., Ms., Dr., Lt Col., etc
    EnumSuffix      //Jr., Sr., III, PhD, Esq, etc
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
    Pics [EnumNumProfileImageFields] *image.Image
    //... others
}

// A port replicator so the usecases can treat the repo's separately
type stringRepo struct {
    Pr *profileRepo
    whichField int
}
func NewStringRepo(pr *profileRepo, which int) (*stringRepo) {
    // Should be an assert
    if which >= EnumNumProfileStringFields {
        log.Fatal(fmt.Errorf("invalid String Repo Enum"))
    }
    return &stringRepo{
        Pr: pr,
        whichField: which,
    }
}

func (sr *stringRepo) createOrUpdate(id uint64, val string) error {
    // Should be an assert
    if sr.whichField >= EnumNumProfileStringFields {
        return fmt.Errorf("profileRepo: Invalid string field")
    }

    pp , ok := sr.Pr.Profiles[id]
    if ok {
        // The profile exists fill it in
        pp.NameAndBios[sr.whichField] = val
    } else {
        // Create the profile and assign it
        pp = PublicProfile{}
    }
    sr.Pr.Profiles[id] = pp
    return nil
}

func (sr *stringRepo) Create(id uint64, val string) error {
    return sr.createOrUpdate(id,val)
}

func (sr *stringRepo) Update(id uint64, val string) error {
    return sr.createOrUpdate(id,val)
}

func (sr *stringRepo) Delete(id uint64) error {
    err := sr.Update(id,"")
    if err != nil {
        return err
    }
    sr.Pr.DeleteNotify(id)
    return nil
}

func (sr *stringRepo) Retrieve(id uint64) (string, error) {
    val, ok := sr.Pr.Profiles[id]
    if ok {
        ret := val.NameAndBios[sr.whichField]
        return ret,nil
    }
    return "", fmt.Errorf("stringRepo id not found")
}

func (sr *stringRepo) RetrieveCount() (int, error) {
    //Note this is max count
    return len(sr.Pr.Profiles), nil
}

func (sr *stringRepo) RetrieveAll() ([]*string, error) {
    ret := make([]*string,len(sr.Pr.Profiles))
    var i int = 0
    for _, pp := range sr.Pr.Profiles {
        val := pp.NameAndBios[sr.whichField]
        if val != "" {
            ret[i] = &val
            i++
        }
    }
    return ret[0:i],nil
}

// Image Repo making interfaces to the profileRepo, to isolate the usecases from the detail
// that the image repo may or may not be in the same storage as the strings in the profile.
// THe below could be a templatized version of the stringrepo, except that it deals with image.Image types
// instead of strings
type imageRepo struct {
    Pr *profileRepo
    whichField int
}

func NewImageRepo(pr *profileRepo, which int) *imageRepo {
    // Should be an assert
    if which >= EnumNumProfileImageFields {
        log.Fatal(fmt.Errorf("invalid image Repo Enum"))
    }
    return &imageRepo{
        Pr: pr,
        whichField: which,
    }
}

func (ir *imageRepo) createOrUpdate(id uint64, val *image.Image) error {
    // Should be an assert
    if ir.whichField >= EnumNumProfileImageFields {
        return fmt.Errorf("profileRepo: Invalid image field")
    }

    pp , ok := ir.Pr.Profiles[id]
    if ok {
        // The profile exists fill it in
        pp.Pics[ir.whichField] = val
    } else {
        // Create the profile and assign it
        pp = PublicProfile{}
    }
    ir.Pr.Profiles[id] = pp
    return nil
}


func (ir *imageRepo) Create(id uint64, val *image.Image) error {
    return ir.createOrUpdate(id,val)
}

func (ir *imageRepo) Update(id uint64, val *image.Image) error {
    return ir.createOrUpdate(id,val)
}

func (ir *imageRepo) Delete(id uint64) error {
    err := ir.Update(id, nil)
    if err != nil {
        return err
    }
    ir.Pr.DeleteNotify(id)
    return nil
}

func (ir *imageRepo) Retrieve(id uint64) (*image.Image, error) {
    val, ok := ir.Pr.Profiles[id]
    if ok {
        ret := val.Pics[ir.whichField]
        return ret,nil
    }
    return nil, fmt.Errorf("imageRepo id not found")
}

func (ir *imageRepo) RetrieveCount() (int, error) {
    //Note this is max count
    return len(ir.Pr.Profiles), nil
}

func (ir *imageRepo) RetrieveAll() ([]*image.Image, error) {
    ret := make([]*image.Image,len(ir.Pr.Profiles))
    var i int = 0
    for _, pp := range ir.Pr.Profiles {
        val := pp.Pics[ir.whichField]
        if val != nil {
            ret[i] = val
            i++
        }
    }
    return ret[0:i],nil
}

