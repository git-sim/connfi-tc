package service

import (
	"github.com/git-sim/tc/app/domain/repo"
)

// Dependency inversion layer to prevent Account from having to know about the repo ifc
// for example the the Account doesn't have to know about how to check if the email is unique
// AccountService struct
type AccountService struct {
	repo repo.AccountRepo
}

// NewAccountService takes in the account repository
func NewAccountService(newrepo repo.AccountRepo) *AccountService {
	return &AccountService{
		repo: newrepo,
	}
}

// AlreadyExists returns if the account exists
func (s *AccountService) AlreadyExists(email string) bool {
	_, err := s.repo.Retrieve(email)
	if err == nil {
		return true
	}
	return false
}
