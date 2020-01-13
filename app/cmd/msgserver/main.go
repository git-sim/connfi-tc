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
	profUsecase := usecase.NewProfileUsecase(dbProfiles)

	//Initialize internal notifications
	usecase.InitSubscribers(accServ, folUsecase, accUsecase, profUsecase, dbPendingMsgs)
	usecase.InitAccounts(accUsecase)

	router := chi.NewRouter()

	//Setup the middleware to be used with router
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(handlers.SetupCORSHandler)
	//router.Use(middleware.Logger)

	//Routes for /accounts
	router.Route("/accounts", func(router chi.Router) {
		router.Post("/", handlers.CreateAccount(accUsecase)) // POST /accounts
		router.Get("/", handlers.GetAccountList(accUsecase)) // GET  /accounts

		// Subrouter1 for /{accountID}
		router.Route("/{accountID}", func(router chi.Router) {
			router.Use(handlers.AccountCtxFunc(accUsecase))                   // Validates the {accountID}
			router.Get("/", handlers.GetAccount(accUsecase))                  // GET  /accounts/1234
			router.Put("/", handlers.PutAccount(accUsecase))                  // PUT  /accounts/1234
			router.Delete("/", handlers.DeleteAccount(accUsecase))            // DEL  /accounts/1234
			router.Get("/folders", handlers.GetFolderList(folUsecase))        // GET /accounts/1234/folders
			router.Get("/folders/{folderID}", handlers.GetFolder(folUsecase)) // GET /accounts/1234/folders/0?limit=10&page=0&sortorder=-1&sort=0

			// Subrouter2 for ./messages
			router.Route("/messages", func(router chi.Router) {
				router.Post("/", handlers.CreateMessage(msgUsecase)) // POST /accounts/1234/messages  {NewMessage}

				// Subrouter 3 for ./{messageID}
				router.Route("/{messageID}", func(router chi.Router) {
					router.Use(handlers.MessageCtxFunc(folUsecase))
					router.Get("/", handlers.GetMessage())                 // GET  /accounts/1234/messages/10
					router.Put("/", handlers.PutMessage(folUsecase))       // PUT  /accounts/1234/messages/10 {viewed:true}
					router.Delete("/", handlers.DeleteMessage(folUsecase)) // DEL  /accounts/1234/messages/10
				})
			})

			// Subrouter2 for profile
			router.Route("/profile", func(router chi.Router) {
				router.Use(handlers.ProfileCtxFunc(profUsecase))
				router.Get("/", handlers.GetProfile(profUsecase))                     // GET /accounts/1234/profile/
				router.Put("/", handlers.PutProfile(profUsecase))                     // PUT /accounts/1234/profile/
				router.Get("/bio", handlers.GetProfileBio(profUsecase))               // GET /accounts/1234/profile/bio
				router.Put("/bio", handlers.PutProfileBio(profUsecase))               // PUT /accounts/1234/profile/bio
				router.Get("/avatar", handlers.GetProfileAvatar(profUsecase))         // GET /accounts/1234/profile/avatar
				router.Put("/avatar", handlers.PutProfileAvatar(profUsecase))         // PUT /accounts/1234/profile/avatar
				router.Get("/background", handlers.GetProfileBackground(profUsecase)) // GET /accounts/1234/profile/background
				router.Put("/background", handlers.PutProfileBackground(profUsecase)) // PUT /accounts/1234/profile/background
			})

		})
	})

	router.Post("/login", handlers.PostLogin(sessionUsecase, accUsecase))
	router.Post("/logout", handlers.PostLogout(sessionUsecase))

	listenString := "0.0.0.0:8080"
	fmt.Println("Listening at ", listenString)

	if err := http.ListenAndServe(listenString, router); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
