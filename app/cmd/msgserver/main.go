package main

import (
	"fmt"
	"log"
	"net/http"
	_ "sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

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

	//	dbMsgs := ram.NewStructRepo()
	dbPendingMsgs := ram.NewStructRepo()
	dbFolders := ram.NewStructRepo()
	//For the folders, pass a func allowing folder repo to be created on demand
	// still isolates the usecase from knowing the specifics of the repo
	folderFactoryFn := func() repo.Generic { return ram.NewGenericRepo() }

	accServ := service.NewAccountService(dbAccounts)
	sessionUsecase := usecase.NewSessionUsecase(nil, accServ)
	accUsecase := usecase.NewAccountUsecase(dbAccounts, sessionUsecase, accServ)
	folUsecase := usecase.NewFoldersUsecase(dbFolders, folderFactoryFn, accServ)
	//	msgUsecase := usecase.NewMsgUsecase(dbMsgs, dbPendingMsgs, folUsecase, accServ)

	profUcs := &handlers.ProfileUsecases{} //A struct to collect up the profile usecases
	profUcs.StrUsecases[handlers.EnumFirstNameUsecase] = usecase.NewProfileStringUsecase(dbFirstNames)
	profUcs.StrUsecases[handlers.EnumLastNameUsecase] = usecase.NewProfileStringUsecase(dbLastNames)
	profUcs.StrUsecases[handlers.EnumBioUsecase] = usecase.NewProfileStringUsecase(dbBios)
	profUcs.ImageUsecases[handlers.EnumAvatarImageUsecase] = usecase.NewProfileImageUsecase(dbAviImgs)
	profUcs.ImageUsecases[handlers.EnumBgImageUsecase] = usecase.NewProfileImageUsecase(dbBgImgs)

	//Initialize internal notifications
	usecase.InitSubscribers(accServ, folUsecase, accUsecase, dbPendingMsgs)
	usecase.InitAccounts(accUsecase)

	router := chi.NewRouter()
	//Setup the middleware to be used with router
	router.Use(middleware.RealIP)
	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(handlers.SetupCORSHandler)

	//Routes for /accounts
	router.Route("/accounts", func(router chi.Router) {
		router.Post("/", handlers.CreateAccount(accUsecase)) // POST /accounts
		router.Get("/", handlers.GetAccountList(accUsecase)) // GET  /accounts

		// Subrouter1 for /{accountID}
		router.Route("/{accountID}", func(router chi.Router) {
			router.Get("/", handlers.GetAccount(accUsecase))       // GET  /accounts/1234
			router.Put("/", handlers.PutAccount(accUsecase))       // PUT  /accounts/1234
			router.Delete("/", handlers.DeleteAccount(accUsecase)) // DEL  /accounts/1234

			/* 			router.Get("/folders", getFoldersList)            // GET /accounts/1234/folders
			   			router.Get("/folders/{folderID}", getFoldersList) // GET /accounts/1234/folders/0?limit=10&page=0&sortorder=-1&sort=0

			   			// Subrouter2 for ./messages
			   			router.Route("/messages", func(router chi.Router) {
			   				router.Post("/", createMessage)              // POST /accounts/1234/messages  {NewMessage}
			   				router.Get("/", getMessages)                 // GET  /accounts/1234/messages?limit=10&offset=0
			   				router.Get("/{messageID}", getOneMessage)    // GET  /accounts/1234/messages/10
			   				router.Put("/{messageID}", putOneMessage)    // PUT  /accounts/1234/messages/10 {viewed:true}
			   				router.Delete("/{messageID}", deleteMessage) // DEL  /accounts/1234/messages/10
			*/
		})

		/* 			// Subrouter2 for ./threads
			router.Route("/threads", func(router chi.Router) {
				router.Get("/", getThreads)                // GET /accounts/1234/threads?limit=10&offset=0
				router.Get("/{threadID}", getOneThread)    // GET /accounts/1234/threads/10
				router.Put("/{threadID}", putOneThread)    // PUT /accounts/1234/threads/10 {mute:true}
				router.Delete("/{threadID}", deleteThread) // DEL /accounts/1234/threads/10
			})

			// Subrouter2 for profile
			router.Route("/profile", func(router chi.Router) {
				router.Get("/", getProfile)                     // GET /accounts/1234/profile
				router.Put("/", putProfile)                     // PUT /accounts/1234/profile
				router.Get("/bio", getProfileBio)               // GET /accounts/1234/profile/bio
				router.Put("/bio", putProfileBio)               // PUT /accounts/1234/profile/bio
				router.Get("/name", getProfileName)             // GET /accounts/1234/profile/name
				router.Put("/name", putProfileName)             // PUT /accounts/1234/profile/name
				router.Get("/avatar", getProfileAvatar)         // GET /accounts/1234/profile/avatar
				router.Put("/avatar", putProfileAvatar)         // PUT /accounts/1234/profile/avatar
				router.Get("/background", getProfileBackground) // GET /accounts/1234/profile/background
				router.Put("/background", putProfileBackground) // PUT /accounts/1234/profile/background
			})
		})
		*/
	})

	//	mux := http.NewServeMux()
	//	mux.Handle("/login", handlers.HandleLogin(sessionUsecase, accUsecase))
	//	mux.Handle("/logout", handlers.HandleLogout(sessionUsecase))
	//	mux.Handle("/account", handlers.HandleAccount(accUsecase))
	//	mux.Handle("/accountList", handlers.HandleAccountList(accUsecase))
	//	mux.Handle("/profile", handlers.HandleProfile(accUsecase, profUcs))
	//	mux.Handle("/message", handlers.HandleMessage(msgUsecase, folUsecase, accUsecase))
	//	mux.Handle("/folder", handlers.HandleFolder(folUsecase, msgUsecase, accUsecase))

	listenString := "0.0.0.0:8080"
	fmt.Println("Listening at ", listenString)

	if err := http.ListenAndServe(listenString, router); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
