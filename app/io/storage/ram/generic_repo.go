package ram

import (
	"fmt"
	"sync"

	"github.com/git-sim/tc/app/domain/repo"
)

// GenericRepo Impl of ram based account repository. Just a map[id]interface{}
type genericRepo struct {
	mtx   *sync.Mutex
	elems map[repo.GenericKeyT]interface{}
}

// NewStructRepo a just an alias
func NewStructRepo() repo.Generic {
	return NewGenericRepo()
}

// NewGenericRepo a repository for holding any type (interface{} in go)
func NewGenericRepo() repo.Generic {
	return &genericRepo{
		mtx:   &sync.Mutex{},
		elems: make(map[repo.GenericKeyT]interface{}),
	}
}

func (gr *genericRepo) createOrUpdate(id repo.GenericKeyT, val interface{}) error {
	gr.elems[id] = val
	return nil
}

func (gr *genericRepo) Create(id repo.GenericKeyT, val interface{}) error {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	return gr.createOrUpdate(id, val)
}

func (gr *genericRepo) Update(id repo.GenericKeyT, val interface{}) error {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	return gr.createOrUpdate(id, val)
}

func (gr *genericRepo) Delete(id repo.GenericKeyT) error {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	delete(gr.elems, id)
	return nil
}

func (gr *genericRepo) Retrieve(id repo.GenericKeyT) (interface{}, error) {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	val, ok := gr.elems[id]
	if ok {
		ret := val
		return ret, nil
	}
	return nil, fmt.Errorf("genericRepo id not found")
}

func (gr *genericRepo) RetrieveCount() (int, error) {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	//Note this is max count
	return len(gr.elems), nil
}

func (gr *genericRepo) RetrieveFiltered(fn func(interface{}) bool) ([]interface{}, error) {
	ret, err := gr.RetrieveAll()
	if err != nil {
		return nil, err
	}
	//out := filter(ret,fn)
	var cursor int
	for _, v := range ret {
		if !fn(v) {
			continue //drop
		}
		ret[cursor] = v
		cursor++
	}
	return ret[:cursor], nil
}

func (gr *genericRepo) RetrieveAll() ([]interface{}, error) {
	gr.mtx.Lock()
	defer gr.mtx.Unlock()
	ret := make([]interface{}, len(gr.elems))
	var i int
	for _, e := range gr.elems {
		if e != nil {
			ret[i] = e
			i++
		}
	}
	return ret[0:i], nil
}
