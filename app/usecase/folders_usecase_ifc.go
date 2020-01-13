package usecase

import (
	"strings"
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

var folderEnum = map[string]int{
	"inbox":     EnumInbox,
	"archive":   EnumArchive,
	"sent":      EnumSent,
	"scheduled": EnumScheduled,
}

// FolderEnum returns empty string if invalid
func FolderEnum(in string) int {
	v, ok := folderEnum[strings.ToLower(in)]
	if !ok {
		return EnumInbox
	}
	return v
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
	FolderIdx int `json:"folderid"`
	SortBy    int `json:"sort"`
	SortOrder int `json:"sortorder"`
	Limit     int `json:"limit"`
	Page      int `json:"page"`
}

type MsgQueryOutput struct {
	Requested   QueryParams
	QueriedAt   time.Time  `json:"queriedat"`
	FolderName  string     `json:"foldername"`
	NumTotal    int        `json:"numtotal"`
	NumUnviewed int        `json:"numunviewed"`
	NumElems    int        `json:"numelems"`
	Elems       []MsgEntry `json:"elems"`
}

type FolderInfoOutput struct {
	FolderIdx  int    `json:"folderid"`
	FolderName string `json:"foldername"`
	NumTotal   int    `json:"numtotal"`
}

// FoldersUsecase handles folder management
type FoldersUsecase interface {

	// Controller related functionality ---
	AddToFolder(folderEnum int, id AccountIDType, msg MsgEntry) error
	UpdateMsg(id AccountIDType, mid MsgIDType, msg MsgEntry) error

	UpdateViewed(id AccountIDType, mid MsgIDType, newval bool) error
	UpdateStarred(id AccountIDType, mid MsgIDType, newval bool) error
	ArchiveMsg(id AccountIDType, mid MsgIDType) error
	UnArchiveMsg(id AccountIDType, mid MsgIDType) error
	DeleteMsg(id AccountIDType, mid MsgIDType) error

	// Presenter Functions ---
	QueryMsgs(id AccountIDType, qp QueryParams) (*MsgQueryOutput, error)
	GetOneMsg(aID AccountIDType, mID MsgIDType) (*MsgEntry, error)

	//List
	GetFolderInfo(id AccountIDType) ([]FolderInfoOutput, error)

	// For use by the system ---
	//CreateNewFolders ... called by NofityNewAccount
	CreateNewFolders(acc entity.Account) error
	//Delete
}
