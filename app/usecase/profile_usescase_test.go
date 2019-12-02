package usecase
import (
    "testing"
    "os"
    "bufio"
    "image"
    _ "image/png"
    
    "reflect"
)

// mockStringRepo ---
// Note for testing only doesn't do any locking 
type stringRepo struct {
    m map[string]string
}

func (r *stringRepo) Create(email string, val string) error {
    r.m[email] = val
    return nil
}
    
func (r *stringRepo) Update(email string, val string) error {
    r.m[email] = val
    return nil
}
func (r *stringRepo) Delete(email string) error {
    delete(r.m, email)
    return nil
}

func (r *stringRepo) Retrieve(email string) (string, error) {
    a, ok := r.m[email]
    if ok {
        return a, nil
    } else {        
        return "", nil
    }
}

func (r *stringRepo) RetrieveCount() (int, error) {
    return len(r.m), nil
}

func (r *stringRepo) RetrieveAll() ([]string, error) {
    as := []string{}
    for _ , v := range r.m {
        as = append(as, v)
    }
    return as, nil
}

//
func TestProfileNameUsecase(t *testing.T) {
    
    //Setup a mockrepo and fill it with some accounts
    firstNames := []string { "Alice", "Bo" }
    lastNames  := []string { "Smith", "bSmith" }
    
    mockRepoF := &stringRepo {
            m: make(map[string]string),
    }
    mockRepoL := &stringRepo {
            m: make(map[string]string),
    }
    
    as := []string { "Alice.Smith@mail.com", "BobSmith@mail.com", }
    for i , v := range as {
        mockRepoF.Create(v, firstNames[i])
        mockRepoL.Create(v, lastNames[i])
    }
    
    uF, _ := NewProfileFirstNameUsecase(mockRepoF)  
    count, _ := uF.GetCount(); 
    if(count != len(mockRepoF.m)) {
        t.Errorf("Count expected %d got %d",len(mockRepoF.m),count)
    }

    uL, _ := NewProfileLastNameUsecase(mockRepoL)
    countL, _ := uL.GetCount(); 
    if(countL != len(mockRepoL.m)) {
        t.Errorf("Count expected %d got %d",len(mockRepoL.m),countL)
    }
    
    for i , v := range as {
        f, _ := uF.Get(v)
        if !reflect.DeepEqual(f,firstNames[i]) {
            t.Error("expected equal images")
        }
        l, _ := uL.Get(v)
        if !reflect.DeepEqual(l,lastNames[i]) {
            t.Error("expected equal images")
        }
    }

    // Swap names
    for i , v := range as {
        err := uF.Set(v,firstNames[count-1-i])
        if err != nil {
            t.Error("Failed to set firstname")
        }
        errL := uL.Set(v,lastNames[count-1-i])
        if errL != nil {
            t.Error("Failed to set lastname")
        }
    }

    for i , v := range as {
        f, _ := uF.Get(v)
        if !reflect.DeepEqual(f,firstNames[count-1-i]) {
            t.Error("expected equal images")
        }
        l, _ := uL.Get(v)
        if !reflect.DeepEqual(l,lastNames[count-1-i]) {
            t.Error("expected equal images")
        }
    }
}



// mockImageRepo
// Note for testing only doesn't do any locking 
type imageRepo struct {
    m map[string]*image.Image
}

func (r *imageRepo) Create(email string, val *image.Image) error {
    r.m[email] = val
    return nil
}
    
func (r *imageRepo) Update(email string, val *image.Image) error {
    r.m[email] = val
    return nil
}
func (r *imageRepo) Delete(email string) error {
    delete(r.m, email)
    return nil
}

func (r *imageRepo) Retrieve(email string) (*image.Image, error) {
    a, ok := r.m[email]
    if ok {
        return a, nil
    } else {        
        return nil, nil
    }
}

func (r *imageRepo) RetrieveCount() (int, error) {
    return len(r.m), nil
}

func (r *imageRepo) RetrieveAll() ([]*image.Image, error) {
    as := []*image.Image{}
    for _ , v := range r.m {
        as = append(as, v)
    }
    return as, nil
}


//
func TestProfileAvatarUsecase(t *testing.T) {
    
    //Setup a mockrepo and fill it with some accounts
    var fnames = []string { "testdata/img1.png", "testdata/img2.png" }
    var images = []*image.Image{}
    for _, fname := range fnames {
        f, err := os.Open(fname)
        if err != nil {
            t.Errorf("Can't open file %s",fname)
        }
        defer f.Close()
        img, _ , err := image.Decode(bufio.NewReader(f))
        if(err != nil) {
            t.Errorf("Error decoding file %s",fname)
        }
        images = append(images, &img)
    }
    
    mockRepo := &imageRepo {
            m: make(map[string]*image.Image),
    }
    
    as := []string { "Alice.Smith@mail.com", "BobSmith@mail.com", }
    for i , v := range as {
        mockRepo.Create(v, images[i])
    }
    
    u, _ := NewProfileAvatarUsecase(mockRepo)
    count, _ := u.GetCount(); 
    if(count != 2) {
        t.Error("Count != 2")
    }
    for i , v := range as {
        ima, _ := u.Get(v)
        if !reflect.DeepEqual(ima,images[i]) {
            t.Error("expected equal images")
        }
    }
    
    for i , v := range as {
        err := u.Set(v,images[count-1-i])
        if(err != nil) {
            t.Error("Error while setting avatar")
        }
    }

    for i , v := range as {
        ima, _ := u.Get(v)
        if !reflect.DeepEqual(ima,images[count-1-i]) {
            t.Error("expected equal images")
        }
    }   
}
