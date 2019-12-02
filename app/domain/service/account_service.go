package service

import (
    "fmt"
    "github.com/git-sim/tc/app/domain/repo"
)

// Dependency inversion layer to prevent Account from having to know about the repo ifc
// for example the the Account doesn't have to know about how to check if the email is unique

type AccountService struct {
    repo repo.AccountRepo
}

func NewAccountService(newrepo repo.AccountRepo) *AccountService {
    return &AccountService{
        repo: newrepo,
    }
}

func (s *AccountService) AlreadyExists(email string) error {
    a, err := s.repo.Retrieve(email)
    if err != nil { 
        return err
    }   
    if(a != nil) {
        return fmt.Errorf("%s exists", email)
    }
    return nil 
}