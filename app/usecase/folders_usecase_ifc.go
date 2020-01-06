package usecase

import (
	"time"

	"github.com/git-sim/tc/app/domain/entity"
)

const (
	EnumInbox = iota
	EnumArchive
	EnumSent
	EnumScheduled
	EnumNumFolders
)

var folderText = map[int]string{
	EnumInbox:     "inbox",
	EnumArchive:   "archive",
	EnumSent:      "sent",
	EnumScheduled: "scheduled",
}

// FolderText returns empty string if invalid
func FolderText(code int) string {
	return folderText[code]
}

const (
	EnumSortByTime = iota
	EnumSortBySubject
	EnumSortBySender
	EnumNumSortBy
)

var sortText = map[int]string{
	EnumSortByTime:    "time",
	EnumSortBySubject: "subject",
	EnumSortBySender:  "sender",
}

//SortText returns empty string if invalid
func SortText(code int) string {
	return sortText[code]
}

type QueryParams struct {
	FolderIdx int
	SortBy    int
	SortOrder int
	Limit     int
	Page      int
}

type MsgQueryOutput struct {
	Requested   QueryParams
	QueriedAt   time.Time
	FolderName  string
	NumTotal    int
	NumUnviewed int
	NumElems    int //in the folder
	Elems       []MsgEntry
}

// FoldersUsecase handles folder management
type FoldersUsecase interface {

	// Controller related functionality
	AddToFolder(folderEnum int, id AccountIDType, msg MsgEntry) error

	UpdateViewed(id AccountIDType, mid MsgIDType, newval bool) error
	UpdateStarred(id AccountIDType, mid MsgIDType, newval bool) error
	ArchiveMsg(id AccountIDType, mid MsgIDType) error
	UnArchiveMsg(id AccountIDType, mid MsgIDType) error
	DeleteMsg(id AccountIDType, mid MsgIDType) error

	// Presenter Functions
	QueryMsgs(id AccountIDType, qp QueryParams) (*MsgQueryOutput, error)
	//QueryThreads()

	//List

	// For use by the system
	//CreateNewFolders ... called by NofityNewAccount
	CreateNewFolders(acc entity.Account) error
	//Delete
}
