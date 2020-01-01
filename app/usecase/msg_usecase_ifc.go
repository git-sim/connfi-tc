package usecase

import (
	"fmt"
	"strconv"

	"github.com/git-sim/tc/app/domain/entity"
)

// MsgUsecase ifc
type MsgUsecase interface {
	// Checks whether the msg is valid before Enqueuing
	IsValid(msg *IngressMsg) (bool, error)

	// Enqueues the message into the system. On success (newMsgid,nil) on fail (0,err)
	EnqueueMsg(msg *IngressMsg) (MsgIDType, error)

	//Get a message from the msg store
	RetrieveMsg(mid MsgIDType) (*EgressMsg, error)
}

// Forward types from the entity layer so clients don't have to unnecessarily reach into the entity directly
// MsgIDType layer of indirection to allow future extension
type MsgIDType entity.MsgIDType

const MsgIDStringBase = entity.MsgIDStringBase
const MsgIDBits = entity.MsgIDBits

// ThreadIDType layer of indirection to allow future extension
type ThreadIDType entity.ThreadIDType

const ThreadIDStringBase = entity.ThreadIDStringBase
const ThreadIDBits = entity.ThreadIDBits

// IngressMsg the type for messages coming into the usecase/interactor layer
type IngressMsg entity.MsgBase

// EgressMsg the type of messages going out of the layer (has user generated metadata attached)
type EgressMsg entity.Msg

// MsgEntry is the type decorated for folders
type MsgEntry entity.MsgEntry

// Conversion functions

// MsgIDToString conversion
func MsgIDToString(id MsgIDType) string {
	return strconv.FormatUint(uint64(id), MsgIDStringBase)
}

// ToMsgID so we're all on the same format
func ToMsgID(msgIDString string) (MsgIDType, error) {
	id, err := strconv.ParseUint(msgIDString, MsgIDStringBase, MsgIDBits)
	if err != nil {
		return MsgIDType(0), NewEs(EsArgConvFail,
			fmt.Sprintf("msgIDString %s", msgIDString))
	}
	return MsgIDType(id), nil
}

//NOTE MsgUsecase just forward declares the entity.Msg types at the usecase boundary.
//  It's not against Clean architecture, since the dependency is still inward, but
//  from an extensibility viewpoint, the better way to do it is have different Msg types
//  to isolate the entity from the usecase boundary.
//  That way they can evolve separately. Something like:
//type IngressMsg struct {
//	ParentMid   entity.MsgIDType
//	ScheduledAt time.Time
//	SenderEmail string
//	Recipients  []string
//	Subject     string
//	Body        []byte
//}

//type EgressMsg struct {
//	Mid      entity.MsgIDType
//	Tid      entity.ThreadIDType
//	SentAt   time.Time
//	SenderID entity.AccountIDType
//	M        IngressMsg
//}
