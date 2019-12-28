package service

import (
	"fmt"
	_ "reflect"
	"testing"

	"github.com/git-sim/tc/app/domain/entity"
)

// mockRepo ---
// Note for testing only doesn't do any locking
type accountRepo struct {
	m map[string]*entity.Account
}

func (r *accountRepo) Create(a *entity.Account) error {
	r.m[a.GetEmail()] = a
	return nil
}

func (r *accountRepo) Update(a *entity.Account) error {
	r.m[a.GetEmail()] = a
	return nil
}
func (r *accountRepo) Delete(a *entity.Account) error {
	delete(r.m, a.GetEmail())
	return nil
}

func (r *accountRepo) Retrieve(email string) (*entity.Account, error) {
	a, ok := r.m[email]
	if ok {
		return a, nil
	} else {
		return nil, fmt.Errorf("email not found")
	}
}

func (r *accountRepo) RetrieveCount() (int, error) {
	return len(r.m), nil
}

func (r *accountRepo) RetrieveAll() ([]*entity.Account, error) {
	as := []*entity.Account{}
	for _, v := range r.m {
		as = append(as, v)
	}
	return as, nil
}

//
func TestAccountService(t *testing.T) {

	//Setup a mockrepo and fill it with some accounts
	mockRepo := &accountRepo{
		m: make(map[string]*entity.Account),
	}

	as := []string{"Alice.Smith@mail.com", "BobSmith@mail.com"}

	bs := entity.NewAccounts(as[0], as[1])
	for _, v := range bs {
		mockRepo.Create(v)
	}

	s := NewAccountService(mockRepo)

	// Test Account service Already Exist with unique account (should pass)
	if s.AlreadyExists("Charlie_Smith@mail.com") == true {
		t.Error("AlreadyExists fails for new email")
	}

	// AlreadyExists with duplicate account should fail
	if s.AlreadyExists(bs[0].GetEmail()) == false {
		t.Error("AlreadyExists failed to catch duplicate email")
	}
}
