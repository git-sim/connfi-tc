package usecase

import (
	"fmt"
	"reflect"
	"sort"
	"time"

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
func (f *foldersUsecase) AddToFolder(folderEnum int, id AccountIDType, msg MsgEntry) error {
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

	enmsg := entity.MsgEntry(msg) //Entity messages go in to the repos
	//Convert to the desired type in folder (unnecessary complexity?)
	// reverted to switches, need to figure out how to switch on ElemT
	switch folderEnum {
	case EnumInbox:
		err = folders[folderEnum].Create(msgkey, enmsg)
	case EnumArchive:
		err = folders[folderEnum].Create(msgkey, enmsg)
	case EnumSent:
		err = folders[folderEnum].Create(msgkey, enmsg.M)
	case EnumScheduled:
		err = folders[folderEnum].Create(msgkey, enmsg.M)
	default:
		err = NewEs(EsInternalError, "Unknown msg type for folder")
	}
	//
	return err
}

func (f *foldersUsecase) ArchiveMsg(id AccountIDType, mid MsgIDType) error {
	return f.moveBetweenFolders(EnumInbox, EnumArchive, id, mid)
}

func (f *foldersUsecase) UnArchiveMsg(id AccountIDType, mid MsgIDType) error {
	return f.moveBetweenFolders(EnumArchive, EnumInbox, id, mid)
}

func (f *foldersUsecase) moveBetweenFolders(srcEnum int, destEnum int, id AccountIDType, mid MsgIDType) error {
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

func (f *foldersUsecase) UpdateStarred(id AccountIDType, mid MsgIDType, newval bool) error {
	return f.updateViewedStarred(id, mid, func(pMsg *entity.MsgEntry) {
		if pMsg != nil {
			pMsg.IsStarred = newval
		}
	})
}

func (f *foldersUsecase) UpdateViewed(id AccountIDType, mid MsgIDType, newval bool) error {
	return f.updateViewedStarred(id, mid, func(pMsg *entity.MsgEntry) {
		if pMsg != nil {
			pMsg.IsViewed = newval
			if newval {
				pMsg.ViewedAt = time.Now()
			} else {
				pMsg.ViewedAt = time.Time{} //empty it out
			}
		}
	})
}

func (f *foldersUsecase) updateViewedStarred(id AccountIDType, mid MsgIDType, fn func(*entity.MsgEntry)) error {

	// Messaging rule: only the messages in the Inbox, Archive (or user folders) have viewed,starred facility
	acckey := repo.GenericKeyT(id)
	val, err := f.dbFolders.Retrieve(acckey)
	if err != nil {
		return err
	}
	folders := val.([EnumNumFolders]repo.Generic)
	msgkey := repo.GenericKeyT(mid)

	folderEnums := []int{EnumInbox, EnumArchive}

	found := false
	for _, folderEnum := range folderEnums {
		val, err := folders[folderEnum].Retrieve(msgkey)
		if err != nil {
			continue //next folder
		}
		found = true

		msg, ok := val.(entity.MsgEntry)
		if ok {
			fn(&msg)
			folders[folderEnum].Update(msgkey, msg)
		}
	}
	if found {
		return nil
	}
	return err
}

// Presenter Funcionality for the Folders
func isValidQuery(qp QueryParams) (bool, error) {
	return true, nil //todo checking already done by the handler, but the boundary needs it's own check
}

// Define out some sorters (these are less functions)
type sorterMsgEntry func(i, j int, reverse bool, coll []MsgEntry) bool

var sortFuncs = map[int]sorterMsgEntry{
	EnumSortByTime: func(i, j int, reverse bool, coll []MsgEntry) bool {
		return coll[i].M.SentAt.Before(coll[j].M.SentAt) != reverse
	},
	EnumSortBySubject: func(i, j int, reverse bool, coll []MsgEntry) bool {
		return (coll[i].M.M.Subject < coll[j].M.M.Subject) != reverse
	},
	EnumSortBySender: func(i, j int, reverse bool, coll []MsgEntry) bool {
		return (coll[i].M.M.SenderEmail < coll[j].M.M.SenderEmail) != reverse
	},
}

func (f *foldersUsecase) QueryMsgs(id AccountIDType, qp QueryParams) (*MsgQueryOutput, error) {
	acckey := repo.GenericKeyT(id)
	val, err := f.dbFolders.Retrieve(acckey)
	if err != nil {
		return nil, err
	}

	if _, err := isValidQuery(qp); err != nil {
		return nil, err
	}

	folders := val.([EnumNumFolders]repo.Generic)
	folder := folders[qp.FolderIdx]

	pOut := &MsgQueryOutput{
		Requested:  qp,
		QueriedAt:  time.Now(),
		FolderName: FolderText(qp.FolderIdx),
	}

	msgs, err := folder.RetrieveAll()
	if err != nil {
		return nil, NewEs(EsInternalError,
			fmt.Sprintf("Unable to retrieve from Folder idx %d", qp.FolderIdx))
	}

	elems := make([]MsgEntry, 0, len(msgs))

	// The folder is a table/collection in the storage repo/database.
	// todo Take advantage of the query capabilities in repositories.
	// For now since the data is opaque to the db, we pull it and give it form here
	// Also need to refactor this to put the type switching/dispatch in one call or interface
	switch qp.FolderIdx {
	case EnumSent:
		fallthrough
	case EnumScheduled:
		// Sent & Scheduled contains type entity.Msg
		for _, val := range msgs {
			val2, _ := val.(entity.Msg)
			elem := MsgEntry(*entity.NewMsgEntry(val2))
			elems = append(elems, elem)
		}

	default:
		// Inbox,Archive and all user folders contain entity.MsgEntry
		numUnviewed := 0
		for _, val := range msgs {
			val2, _ := val.(entity.MsgEntry)
			elem := MsgEntry(val2)
			elems = append(elems, elem)
			if !val2.IsViewed {
				numUnviewed++
			}
		}
		pOut.NumUnviewed = numUnviewed
	}

	nelems := len(elems)
	pOut.NumTotal = nelems
	// Ok we have the data in the elems array now sort it, select it out
	sort.Slice(elems, func(i, j int) bool {
		return sorterMsgEntry(sortFuncs[qp.SortBy])(i, j, qp.SortOrder == 0, elems)
	})
	// Offset and trim the response, make sure we never send more than the limit
	startIdx := qp.Page * qp.Limit
	if startIdx >= nelems {
		// past the end no items to add
		pOut.NumElems = 0
		return pOut, nil
	}
	nToSend := qp.Limit //min(limit,nelems-startIdx)
	if nelems-startIdx < nToSend {
		nToSend = nelems - startIdx
	}

	pOut.NumElems = nToSend
	if nToSend > 0 {
		pOut.Elems = elems[startIdx : startIdx+nToSend] //[a:a] is 0 len slice in go which works here
	} else {
		// bad etiquette to send null, send empty list instead
		pOut.Elems = []MsgEntry{}
	}

	return pOut, nil
}
