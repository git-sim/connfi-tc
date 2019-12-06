package main

import (
	"fmt"
    _ "sync"
    _ "github.com/git-sim/tc/app/domain/entity"
    _ "github.com/git-sim/tc/app/usecase"
    _ "github.com/git-sim/tc/app/io/storage/ram"
    _ "github.com/git-sim/tc/app/io/rest"
	"net/http"
	
)

func main() {

	fmt.Println("Hello from msgserver main()")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}