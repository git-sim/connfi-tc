package main

import (
	"fmt"
	"log"
	"net/http"
	_ "sync"

	"github.com/git-sim/tc/app/domain/service"
	"github.com/git-sim/tc/app/io/rest/handlers"
	"github.com/git-sim/tc/app/io/storage/ram"
	"github.com/git-sim/tc/app/usecase"
)

func main() {

	fmt.Println("Hello from msgserver main()")
	dbAccounts := ram.NewAccountRepo()
	dbProfiles := ram.NewProfileRepo()
	dbFirstNames := ram.NewStringRepo(dbProfiles, ram.EnumFirstName)
	dbLastNames := ram.NewStringRepo(dbProfiles, ram.EnumLastName)
	dbBios := ram.NewStringRepo(dbProfiles, ram.EnumBio)
	dbAviImgs := ram.NewImageRepo(dbProfiles, ram.EnumAvatar)
	dbBgImgs := ram.NewImageRepo(dbProfiles, ram.EnumBackground)
	dbMsgs := ram.NewStructRepo()
	dbPendingMsgs := ram.NewStructRepo()

	accServ := service.NewAccountService(dbAccounts)
	sessionUsecase := usecase.NewSessionUsecase(nil, accServ)
	accUsecase := usecase.NewAccountUsecase(dbAccounts, sessionUsecase, accServ)
	msgUsecase := usecase.NewMsgUsecase(dbMsgs, dbPendingMsgs, accServ)

	profUcs := &handlers.ProfileUsecases{} //A struct to collect up the profile usecases
	profUcs.StrUsecases[handlers.EnumFirstNameUsecase] = usecase.NewProfileStringUsecase(dbFirstNames)
	profUcs.StrUsecases[handlers.EnumLastNameUsecase] = usecase.NewProfileStringUsecase(dbLastNames)
	profUcs.StrUsecases[handlers.EnumBioUsecase] = usecase.NewProfileStringUsecase(dbBios)
	profUcs.ImageUsecases[handlers.EnumAvatarImageUsecase] = usecase.NewProfileImageUsecase(dbAviImgs)
	profUcs.ImageUsecases[handlers.EnumBgImageUsecase] = usecase.NewProfileImageUsecase(dbBgImgs)

	mux := http.NewServeMux()
	mux.Handle("/login", handlers.HandleLogin(sessionUsecase, accUsecase))
	mux.Handle("/logout", handlers.HandleLogout(sessionUsecase))
	mux.Handle("/account", handlers.HandleAccount(accUsecase))
	mux.Handle("/accountList", handlers.HandleAccountList(accUsecase))
	mux.Handle("/profile", handlers.HandleProfile(accUsecase, profUcs))
	mux.Handle("/message", handlers.HandleMessage(msgUsecase, accUsecase))
	//mux.Handle("/directory",  handlers.HandleDirectory(dirUsecase))

	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
