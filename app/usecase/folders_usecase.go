package usecase

import (
	"fmt"
	"reflect"

	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
)

type foldersUsecase struct {
	dbFolders   repo.Generic        // This is a container of collections map[accId][]Folder
	dbFactoryFn func() repo.Generic // Function used to instantiate new collections
	service     *service.AccountService
}

// Define what folders the user starts with
// InboxFolderType def
type InboxFolderType map[entity.MsgIDType]entity.MsgEntry
type ArchiveFolderType map[entity.MsgIDType]entity.MsgEntry
type SentFolderType map[entity.MsgIDType]entity.Msg
type ScheduledFolderType map[entity.MsgIDType]entity.MsgBase

// Some helper structs for table driven code make it easier to add remove folders in the design
type folderDesc struct {
	Name  string
	T     reflect.Type
	ElemT reflect.Type
}

var folderDescs = [EnumNumFolders]folderDesc{
	folderDesc{"Inbox", reflect.TypeOf(InboxFolderType{}), reflect.TypeOf(entity.MsgEntry{})},
	folderDesc{"Archive", reflect.TypeOf(ArchiveFolderType{}), reflect.TypeOf(entity.MsgEntry{})},
	folderDesc{"Sent", reflect.TypeOf(SentFolderType{}), reflect.TypeOf(entity.Msg{})},
	folderDesc{"Scheduled", reflect.TypeOf(ScheduledFolderType{}), reflect.TypeOf(entity.MsgBase{})},
}

// NewFoldersUsecase ctor
func NewFoldersUsecase(dbFolders repo.Generic, dbFactoryFn func() repo.Generic, service *service.AccountService) FoldersUsecase {
	// Create a base repository
	return &foldersUsecase{
		dbFolders:   dbFolders,
		dbFactoryFn: dbFactoryFn,
		service:     service,
	}
}

// CreateNewFolders for a new account
func (f *foldersUsecase) CreateNewFolders(acc entity.Account) error {

	//Create the array of folders for each account
	var folderArray = [EnumNumFolders]repo.Generic{}
	for i := 0; i < EnumNumFolders; i++ {
		folderArray[i] = f.dbFactoryFn()
	}
	// Add it to the dbFolders
	f.dbFolders.Create(repo.GenericKeyT(acc.GetID()), folderArray)
	return nil
}

// Add a message to a user's folder
func (f *foldersUsecase) AddToFolder(folderEnum int, id entity.AccountIDType, msg entity.MsgEntry) error {
	// should be an assert
	if folderEnum < 0 || EnumNumFolders <= folderEnum {
		return NewEs(EsArgInvalid,
			fmt.Sprintf("folderEnum %d", folderEnum))
	}

	acckey := repo.GenericKeyT(id)
	val, err := f.dbFolders.Retrieve(acckey)
	if err != nil {
		return err
	}
	folders := val.([EnumNumFolders]repo.Generic)
	msgkey := repo.GenericKeyT(msg.Mid)
	//Convert to the desired type in folder (unnecessary complexity?)
	// reverted to switches, need to figure out how to switch on ElemT
	switch folderEnum {
	case EnumInbox:
		err = folders[folderEnum].Create(msgkey, msg)
	case EnumArchive:
		err = folders[folderEnum].Create(msgkey, msg)
	case EnumSent:
		err = folders[folderEnum].Create(msgkey, msg.M)
	case EnumScheduled:
		err = folders[folderEnum].Create(msgkey, msg.M.M)
	default:
		err = NewEs(EsInternalError, "Unknown msg type for folder")
	}
	//
	return err
}

// Original that was working
// // Add a message to a user's folder
// func (f *foldersUsecase) AddToFolder(folderEnum int, id entity.AccountIDType, msg entity.MsgEntry) error {
// 	// should be an assert
// 	if folderEnum < 0 || EnumNumFolders <= folderEnum {
// 		return NewEs(EsArgInvalid,
// 			fmt.Sprintf("folderEnum %d", folderEnum))
// 	}

// 	// todo fix this. It is a hack because I haven't looked up the
// 	// best way to do a collection of a collection, preserving the interface.
// 	// Right now I get the whole folder, update it and put it back. This is
// 	// a glaring concurrency problem.  Come back to this after the frontend ifc if worked out
// 	dbkey := repo.GenericKeyT(id)
// 	val, err := f.dbFolders.Retrieve(dbkey)
// 	if err != nil {
// 		return err
// 	}
// 	inbox := val.(InboxFolderType)
// 	inbox[msg.Mid] = msg
// 	f.dbFolders.Update(dbkey, inbox)
// 	//
// 	return nil
// }

func (f *foldersUsecase) ArchiveMsg(id entity.AccountIDType, mid entity.MsgIDType) error {
	return f.moveBetweenFolders(EnumInbox, EnumArchive, id, mid)
}

func (f *foldersUsecase) moveBetweenFolders(srcEnum int, destEnum int, id entity.AccountIDType, mid entity.MsgIDType) error {
	// should be an assert
	if srcEnum < 0 || EnumNumFolders <= srcEnum {
		return NewEs(EsArgInvalid,
			fmt.Sprintf("srcEnum %d", srcEnum))
	}

	if destEnum < 0 || EnumNumFolders <= destEnum {
		return NewEs(EsArgInvalid,
			fmt.Sprintf("destEnum %d", destEnum))
	}

	if srcEnum == destEnum {
		return nil
	}

	// Todo move.  Get msg from src move to dest
	return nil
}
