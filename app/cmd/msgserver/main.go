package main

import (
    "fmt"
    "log"
    _ "sync"
    "net/http"
    "github.com/git-sim/tc/app/domain/service"
    "github.com/git-sim/tc/app/io/storage/ram"
    "github.com/git-sim/tc/app/io/rest/handlers"
    "github.com/git-sim/tc/app/usecase"
)

func main() {

    fmt.Println("Hello from msgserver main()")
	db         := ram.NewAccountRepo()
	accServ    := service.NewAccountService(db)
	accUsecase := usecase.NewAccountUsecase(db, accServ)
	
	mux := http.NewServeMux()
	mux.Handle("/account",    handlers.HandleAccount(accUsecase))
	mux.Handle("/accountList",handlers.HandleAccountList(accUsecase))
	
    if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
        log.Fatal("ListenAndServer:", err)
    }
}
