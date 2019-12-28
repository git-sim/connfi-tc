package usecase

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/git-sim/tc/app/domain/entity"
	"github.com/git-sim/tc/app/domain/repo"
	"github.com/git-sim/tc/app/domain/service"
)

type msgUsecase struct {
	dbMsg     repo.Generic
	dbPending repo.Generic
	service   *service.AccountService
}

// Msg and Thread Id helpers
var lastMsgID uint64
var lastThreadID uint64

func getNextMsgID() MsgIDType {
	return MsgIDType(atomic.AddUint64(&lastMsgID, 1))
}

func getNextThreadID() ThreadIDType {
	return ThreadIDType(atomic.AddUint64(&lastThreadID, 1))
}

// NewMsgUsecase news usecase
func NewMsgUsecase(dbMsg repo.Generic, dbPending repo.Generic, service *service.AccountService) MsgUsecase {
	return &msgUsecase{
		dbMsg:     dbMsg,
		dbPending: dbPending,
		service:   service,
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
	Err error
}

func (ce *checkErr) Check(newErrStr string, assertFn func() bool) {
	if ce.Err != nil {
		return
	}
	if assertFn() == false {
		ce.Err = errors.New(newErrStr)
	}
	return
}

// EnqueueMsg adds the message to the system for scheduling/delivery
func (u *msgUsecase) IsValid(msg *IngressMsg) (bool, error) {
	// The follwing series of checks executed while err stays nil
	ce := &checkErr{}
	ce.Check("Invalid SenderEmail format",
		func() bool { return IsValidEmailStr(msg.SenderEmail) })

	ce.Check("No Recipients",
		func() bool { return len(msg.Recipients) > 0 })

	ce.Check("No Valid Recipient email formats",
		func() bool {
			for _, rcp := range msg.Recipients {
				if IsValidEmailStr(rcp) {
					return true
				}
			}
			return false
		})

	ce.Check("Sender is not registered",
		func() bool {
			return u.service.AlreadyExists(msg.SenderEmail)
		})

	var ok = (ce.Err == nil)
	return ok, ce.Err
}

// EnqueueMsg adds the message to the system for scheduling/delivery
func (u *msgUsecase) EnqueueMsg(msg *IngressMsg) (*EgressMsg, error) {
	if ok, err := u.IsValid(msg); !ok {
		return nil, err
	}

	newmsg := entity.Msg{M: entity.MsgBase(*msg)}

	//Validate or Assign ThreadId
	if msg.ParentMid == 0 {
		newmsg.Tid = entity.ThreadIDType(getNextThreadID())
	} else {
		//todo find the parent message, and assign it's thread id

	}
	//Validate ParentMsgId
	//
	//Fill in SenderID
	//Assign new MsgId
	newmsg.Mid = entity.MsgIDType(getNextMsgID())
	//Add to message store
	err := u.dbMsg.Create(repo.GenericKeyT(newmsg.Mid), newmsg)
	if err != nil {

	}

	//Check if Scheduled for future delivery
	//if so add to sender's Scheduled folder and Add timer
	if msg.ScheduledAt.After(time.Now().Add(time.Second * 10)) {
		//todo
	} else {
		//else assign SentAt and Dispatch the message to recipients
		newmsg.SentAt = time.Now()
		// todo dispatch
	}

	//Check if all Recipients existed, if not add to PendingMsgs with the missing recips
	//
	return nil, nil
}

// RetrieveMsg gets the specified message
func (u *msgUsecase) RetrieveMsg(mid MsgIDType) (*EgressMsg, error) {
	return nil, nil
}
