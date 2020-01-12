package usecase

import (
	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
)

// Register, connect up subscribers for the events in the system

// InitAccounts ...
func InitAccounts(accUsecase AccountUsecase) error {
	_, err := accUsecase.RegisterAccountByEmail("admin@localhost")
	return err
}

// InitSubscribers called at bootup
func InitSubscribers(accServ *service.AccountService, folUsecase FoldersUsecase,
	accUsecase AccountUsecase, dbPendingMsgs repo.Generic) error {
	err := initRegisterAccountSubsribers(accServ, folUsecase, accUsecase, dbPendingMsgs)
	if err != nil {
		return err
	}
	err = initDeleteAccountSubsribers(accServ, folUsecase, accUsecase)
	if err != nil {
		return err
	}
	err = initEnqueueMsgSubsribers(accServ, folUsecase, accUsecase)
	if err != nil {
		return err
	}

	return err
}

func initRegisterAccountSubsribers(accServ *service.AccountService, folUsecase FoldersUsecase,
	accUsecase AccountUsecase, dbPendingMsgs repo.Generic) error {

	accServ.SubscribeRegisterAccount(
		func(acc entity.Account) {
			// Create the new folders for the user
			folUsecase.CreateNewFolders(acc)
		})

	accServ.SubscribeRegisterAccount(
		// Scan pending messages looking for any meant for the newly created recipient
		func(acc entity.Account) {
			// Ugly but it works, there's no concurrency issue because the PendingMsg has been
			// duplicated for each missing recipient so they'll only update their copy.
			pendArray, err := dbPendingMsgs.RetrieveAll()
			if err == nil {
				for _, val := range pendArray {
					if pendmsg, ok := val.(entity.PendingMsgEntry); ok {
						if pendmsg.RecipientLeft == acc.GetEmail() {
							folUsecase.AddToFolder(EnumInbox,
								AccountIDType(acc.GetID()),
								MsgEntry(pendmsg.E))

							// Update/Delete the pending msg
							msgkey := repo.GenericKeyT(pendmsg.E.Mid)
							dbPendingMsgs.Delete(msgkey)
						}
					}
				}
			}
		})

	return nil
}

func initDeleteAccountSubsribers(accServ *service.AccountService, folUsecase FoldersUsecase,
	accUsecase AccountUsecase) error {
	return nil //tbd
}
func initEnqueueMsgSubsribers(accServ *service.AccountService, folUsecase FoldersUsecase,
	accUsecase AccountUsecase) error {
	return nil //tbd
}
