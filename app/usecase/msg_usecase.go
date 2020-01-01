package usecase

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
)

type msgUsecase struct {
	dbMsg      repo.Generic
	dbPending  repo.Generic
	folUsecase FoldersUsecase
	service    *service.AccountService
}

// Msg and Thread Id helpers
var lastMsgID uint64
var lastThreadID uint64

// Returns a new unique id, the implementation is a detail
func getNewMsgID() MsgIDType {
	return MsgIDType(atomic.AddUint64(&lastMsgID, 1))
}

func getNewThreadID() ThreadIDType {
	return ThreadIDType(atomic.AddUint64(&lastThreadID, 1))
}

// NewMsgUsecase news usecase
func NewMsgUsecase(dbMsg repo.Generic, dbPending repo.Generic, folUsecase FoldersUsecase, service *service.AccountService) MsgUsecase {
	return &msgUsecase{
		dbMsg:      dbMsg,
		dbPending:  dbPending,
		folUsecase: folUsecase,
		service:    service,
	}
}

// IsValidEmailStr basic checks on the string not that its registered
func IsValidEmailStr(email string) bool {
	ok := true
	ok = ok && len(email) > 0
	ok = ok && len(email) < 80 //todo make this a config/constant
	return ok
}

type checkErr struct {
	Err *ErrStat
}

func (ce *checkErr) Check(newErrCode int, newErrStr string, assertFn func() bool) {
	if ce.Err != nil {
		return
	}
	if assertFn() == false {
		ce.Err = NewEs(newErrCode, newErrStr)
	}
	return
}

// EnqueueMsg adds the message to the system for scheduling/delivery
func (u *msgUsecase) IsValid(msg *IngressMsg) (bool, error) {
	// The follwing series of checks executed while err stays nil
	ce := &checkErr{}
	ce.Check(EsArgInvalid, "SenderEmail format",
		func() bool { return IsValidEmailStr(msg.SenderEmail) })

	ce.Check(EsArgInvalid, "No Recipients",
		func() bool { return len(msg.Recipients) > 0 })

	ce.Check(EsArgInvalid, "No Valid Recipient email formats",
		func() bool {
			for _, rcp := range msg.Recipients {
				if IsValidEmailStr(rcp) {
					return true
				}
			}
			return false
		})

	ce.Check(EsNotFound, "Sender is not registered",
		func() bool {
			return u.service.AlreadyExists(msg.SenderEmail)
		})

	if ce.Err != nil {
		return false, ce.Err
	}
	return true, nil
}

// EnqueueMsg adds the message to the system for scheduling/delivery
func (u *msgUsecase) EnqueueMsg(msg *IngressMsg) (MsgIDType, error) {
	// Sanity check
	if ok, err := u.IsValid(msg); !ok {
		return 0, err
	}

	// Prepare the message struct adding meta data as needed
	//
	newmsg := entity.Msg{M: entity.MsgBase(*msg)}

	//Validate or Assign ThreadId
	if msg.ParentMid == 0 {
		newmsg.Tid = entity.ThreadIDType(getNewThreadID())
	} else {
		//todo find the parent message, and assign it's thread id

	}
	//Validate ParentMsgId //todo thread handling
	//Fill in SenderID
	if senderID, err := u.service.GetIDFromEmail(msg.SenderEmail); err == nil {
		newmsg.SenderID = senderID
	} else {
		// error the sender issue with sender id
		return 0, NewEs(EsNotFound, "Sender Account ID")
	}

	// Assign new MsgId and Store the Message
	//
	newid := getNewMsgID()
	newmsg.Mid = entity.MsgIDType(newid)
	if err := u.dbMsg.Create(repo.GenericKeyT(newmsg.Mid), newmsg); err != nil {
		return 0, err
	}

	// Handle the scheduling and dispatch if needed
	//
	//Check if Scheduled for future delivery
	//if so add to sender's Scheduled folder and Add timer
	if msg.ScheduledAt.After(time.Now().Add(time.Second * 10)) {
		//Put the message in the scheduled folder, notify a timer
		pMsgEntry := entity.NewMsgEntry(newmsg)
		err := u.folUsecase.AddToFolder(EnumScheduled,
			AccountIDType(newmsg.SenderID), MsgEntry(*pMsgEntry))
		if err != nil {
			return newid, NewEs(EsInternalError,
				fmt.Sprintf("Couldn't add to scheduled folder sender %d, %s",
					newmsg.SenderID, err.Error()))
		}

	} else {
		//else assign SentAt and Dispatch the message to recipients
		newmsg.SentAt = time.Now()
		// todo dispatch move to its own function in the entity layer

		// Add to Sent folder
		pMsgEntry := entity.NewMsgEntry(newmsg)
		err := u.folUsecase.AddToFolder(EnumSent, AccountIDType(newmsg.SenderID), MsgEntry(*pMsgEntry))
		if err != nil {
			return newid, NewEs(EsInternalError,
				fmt.Sprintf("%s", err.Error()))
		}
		// Dispatch to recipients
		for _, recip := range newmsg.M.Recipients {
			recipID, err := u.service.GetIDFromEmail(recip)
			if err == nil {
				// Recipient is in the system send message
				err = u.folUsecase.AddToFolder(EnumInbox, AccountIDType(recipID), MsgEntry(*pMsgEntry))
				if err != nil {
					return newid, NewEs(EsInternalError,
						fmt.Sprintf("%s", err.Error()))
				}
			} else {
				// Recipient isn't in the system, add the message to the pending queue
				pNewPendMsg := entity.NewPendingMsgEntry(*pMsgEntry, recip)
				err = u.dbPending.Create(repo.GenericKeyT(pNewPendMsg.E.Mid), *pNewPendMsg)
				if err != nil {
					return newid, NewEs(EsInternalError,
						fmt.Sprintf("%s", err.Error()))
				}
			}
		}
	}
	return newid, nil
}

// RetrieveMsg gets the specified message from the message store
func (u *msgUsecase) RetrieveMsg(mid MsgIDType) (*EgressMsg, error) {
	val, err := u.dbMsg.Retrieve(repo.GenericKeyT(mid))
	if err != nil {
		return nil, NewEs(EsNotFound,
			fmt.Sprintf("Message with id %d", mid))
	}
	valAsEnt, ok := val.(entity.Msg) //Type assert first
	if !ok {
		return nil, NewEs(EsArgConvFail, "Repository to entity.Msg")
	}

	emsg := EgressMsg(valAsEnt) //convert to outgoing type
	return &emsg, nil
}
