package usecase

import "github.com/git-sim/tc/app/domain/entity"

const (
	EnumInbox = iota
	EnumArchive
	EnumSent
	EnumScheduled
	EnumNumFolders
)

// FoldersUsecase handles folder management
type FoldersUsecase interface {
	AddToFolder(folderEnum int, id entity.AccountIDType, msg entity.MsgEntry) error
	ArchiveMsg(id entity.AccountIDType, mid entity.MsgIDType) error
	//Create
	//Create
	CreateNewFolders(acc entity.Account) error
	//Delete
	//List
}
