package usecase

import (
	"bufio"
	"image"
	_ "image/png"
	"os"
	"testing"

	"reflect"
)

// mockStringRepo ---
// Note for testing only doesn't do any locking
type stringRepo struct {
	m map[uint64]*string
}

func (r *stringRepo) Create(id uint64, val string) error {
	r.m[id] = &val
	return nil
}

func (r *stringRepo) Update(id uint64, val string) error {
	r.m[id] = &val
	return nil
}
func (r *stringRepo) Delete(id uint64) error {
	delete(r.m, id)
	return nil
}

func (r *stringRepo) Retrieve(id uint64) (string, error) {
	a, ok := r.m[id]
	if ok {
		return *a, nil
	}
	return "", nil
}

func (r *stringRepo) RetrieveCount() (int, error) {
	return len(r.m), nil
}

func (r *stringRepo) RetrieveAll() ([]*string, error) {
	as := []*string{}
	for _, v := range r.m {
		as = append(as, v)
	}
	return as, nil
}

//
func TestProfileNameUsecase(t *testing.T) {

	//Setup a mockrepo and fill it with some accounts
	firstNames := []string{"Alice", "Bo"}
	lastNames := []string{"Smith", "bSmith"}

	mockRepoF := &stringRepo{
		m: make(map[uint64]*string),
	}
	mockRepoL := &stringRepo{
		m: make(map[uint64]*string),
	}

	as := []string{"Alice.Smith@mail.com", "BobSmith@mail.com"}
	for i, _ := range as {
		mockRepoF.Create(uint64(i), firstNames[i])
		mockRepoL.Create(uint64(i), lastNames[i])
	}

	uF := NewProfileStringUsecase(mockRepoF)
	count, _ := uF.GetCount()
	if count != len(mockRepoF.m) {
		t.Errorf("count expected %d got %d", len(mockRepoF.m), count)
	}

	uL := NewProfileStringUsecase(mockRepoL)
	countL, _ := uL.GetCount()
	if countL != len(mockRepoL.m) {
		t.Errorf("count expected %d got %d", len(mockRepoL.m), countL)
	}

	for i, _ := range as {
		f, _ := uF.Get(uint64(i))
		if !reflect.DeepEqual(f, firstNames[i]) {
			t.Error("expected equal images")
		}
		l, _ := uL.Get(uint64(i))
		if !reflect.DeepEqual(l, lastNames[i]) {
			t.Error("expected equal images")
		}
	}

	// Swap names
	for i, _ := range as {
		err := uF.Set(uint64(i), firstNames[count-1-i])
		if err != nil {
			t.Error("Failed to set firstname")
		}
		errL := uL.Set(uint64(i), lastNames[count-1-i])
		if errL != nil {
			t.Error("Failed to set lastname")
		}
	}

	for i, _ := range as {
		f, _ := uF.Get(uint64(i))
		if !reflect.DeepEqual(f, firstNames[count-1-i]) {
			t.Error("expected equal images")
		}
		l, _ := uL.Get(uint64(i))
		if !reflect.DeepEqual(l, lastNames[count-1-i]) {
			t.Error("expected equal images")
		}
	}
}

// mockImageRepo
// Note for testing only doesn't do any locking
type imageRepo struct {
	m map[uint64]*image.Image
}

func (r *imageRepo) Create(id uint64, val *image.Image) error {
	r.m[id] = val
	return nil
}

func (r *imageRepo) Update(id uint64, val *image.Image) error {
	r.m[id] = val
	return nil
}
func (r *imageRepo) Delete(id uint64) error {
	delete(r.m, id)
	return nil
}

func (r *imageRepo) Retrieve(id uint64) (*image.Image, error) {
	a, ok := r.m[id]
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
	for _, v := range r.m {
		as = append(as, v)
	}
	return as, nil
}

//
func TestProfileAvatarUsecase(t *testing.T) {

	//Setup a mockrepo and fill it with some accounts
	var fnames = []string{"testdata/img1.png", "testdata/img2.png"}
	var images = []*image.Image{}
	for _, fname := range fnames {
		f, err := os.Open(fname)
		if err != nil {
			t.Errorf("Can't open file %s", fname)
		}
		defer f.Close()
		img, _, err := image.Decode(bufio.NewReader(f))
		if err != nil {
			t.Errorf("Error decoding file %s", fname)
		}
		images = append(images, &img)
	}

	mockRepo := &imageRepo{
		m: make(map[uint64]*image.Image),
	}

	as := []string{"Alice.Smith@mail.com", "BobSmith@mail.com"}
	for i, _ := range as {
		mockRepo.Create(uint64(i), images[i])
	}

	u := NewProfileImageUsecase(mockRepo)
	count, _ := u.GetCount()
	if count != 2 {
		t.Error("Count != 2")
	}
	for i, _ := range as {
		ima, _ := u.Get(uint64(i))
		if !reflect.DeepEqual(ima, images[i]) {
			t.Error("expected equal images")
		}
	}

	for i, _ := range as {
		err := u.Set(uint64(i), images[count-1-i])
		if err != nil {
			t.Error("Error while setting avatar")
		}
	}

	for i, _ := range as {
		ima, _ := u.Get(uint64(i))
		if !reflect.DeepEqual(ima, images[count-1-i]) {
			t.Error("expected equal images")
		}
	}
}
