package entity

import (
    "testing"
    "reflect"
)

func TestNewAccount(t *testing.T) {
    a := &Account{
        id:        1,
        email:     "Alice.Smith@mail.com",
    }
    if !reflect.DeepEqual(a, NewAccount(a.id, a.email)) {
        t.Error("expected equal structs")
    }
}

func TestNewAccounts(t *testing.T) {
    a := []*Account{
        {
            id:        1,
            email:     "Alice.Smith@mail.com",
        },
        {
            id:        2,
            email:     "BobSmith@mail.com",
        },
    }
    if !reflect.DeepEqual(a, NewAccounts(a[0].email, a[1].email)) {
        t.Error("expected equal structs")
    }
}