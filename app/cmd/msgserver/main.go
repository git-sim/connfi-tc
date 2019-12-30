package main

import (
	"fmt"
	"log"
	"net/http"
	_ "sync"

	"github.com/git-sim/tc/app/domain/repo"
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
	dbFolders := ram.NewStructRepo()
	//For the folders, pass a func allowing folder repo to be created on demand
	// still isolates the usecase from knowing the specifics of the repo
	folderFactoryFn := func() repo.Generic { return ram.NewGenericRepo() }

	accServ := service.NewAccountService(dbAccounts)
	sessionUsecase := usecase.NewSessionUsecase(nil, accServ)
	accUsecase := usecase.NewAccountUsecase(dbAccounts, sessionUsecase, accServ)
	folUsecase := usecase.NewFoldersUsecase(dbFolders, folderFactoryFn, accServ)
	msgUsecase := usecase.NewMsgUsecase(dbMsgs, dbPendingMsgs, folUsecase, accServ)

	profUcs := &handlers.ProfileUsecases{} //A struct to collect up the profile usecases
	profUcs.StrUsecases[handlers.EnumFirstNameUsecase] = usecase.NewProfileStringUsecase(dbFirstNames)
	profUcs.StrUsecases[handlers.EnumLastNameUsecase] = usecase.NewProfileStringUsecase(dbLastNames)
	profUcs.StrUsecases[handlers.EnumBioUsecase] = usecase.NewProfileStringUsecase(dbBios)
	profUcs.ImageUsecases[handlers.EnumAvatarImageUsecase] = usecase.NewProfileImageUsecase(dbAviImgs)
	profUcs.ImageUsecases[handlers.EnumBgImageUsecase] = usecase.NewProfileImageUsecase(dbBgImgs)

	//Initialize internal notifications
	usecase.InitSubscribers(accServ, folUsecase, accUsecase, dbPendingMsgs)

	mux := http.NewServeMux()
	mux.Handle("/login", handlers.HandleLogin(sessionUsecase, accUsecase))
	mux.Handle("/logout", handlers.HandleLogout(sessionUsecase))
	mux.Handle("/account", handlers.HandleAccount(accUsecase))
	mux.Handle("/accountList", handlers.HandleAccountList(accUsecase))
	mux.Handle("/profile", handlers.HandleProfile(accUsecase, profUcs))
	mux.Handle("/message", handlers.HandleMessage(msgUsecase, accUsecase))
	//mux.Handle("/directory",  handlers.HandleDirectory(dirUsecase))

	listenString := "0.0.0.0:8080"
	fmt.Println("Listening at ", listenString)

	if err := http.ListenAndServe(listenString, mux); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
