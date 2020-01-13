package ram

import (
	_ "image/png"
	"sync"

	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/usecase"
)

// Impl of ram based profile repository
type profileRepo struct {
	mtx      *sync.Mutex
	Profiles map[uint64]repo.PublicProfile
}

func NewProfileRepo() *profileRepo {
	return &profileRepo{
		mtx:      &sync.Mutex{},
		Profiles: map[uint64]repo.PublicProfile{},
	}
}

func (pr *profileRepo) createOrUpdate(id uint64, val *repo.PublicProfile) error {
	if val == nil {
		return usecase.NewEs(usecase.EsArgInvalid, "*repo.PublicProfile")
	}

	pr.Profiles[id] = *val
	return nil
}

func (pr *profileRepo) Create(id uint64, val *repo.PublicProfile) error {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	return pr.createOrUpdate(id, val)
}

func (pr *profileRepo) Update(id uint64, val *repo.PublicProfile) error {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	return pr.createOrUpdate(id, val)
}

func (pr *profileRepo) Delete(id uint64) error {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	delete(pr.Profiles, id)
	return nil
}

func (pr *profileRepo) Retrieve(id uint64) (*repo.PublicProfile, error) {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	val, ok := pr.Profiles[id]
	if !ok {
		return nil, usecase.NewEs(usecase.EsNotFound, "entity.PublicProfile")
	}
	return &val, nil
}

func (pr *profileRepo) RetrieveCount() (int, error) {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	return len(pr.Profiles), nil
}
