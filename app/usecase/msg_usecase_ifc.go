package usecase

import (
	"github.com/git-sim/tc/app/domain/entity"
)

// MsgIDType layer of indirection to allow future extension
type MsgIDType entity.MsgIDType

// ThreadIDType layer of indirection to allow future extension
type ThreadIDType entity.ThreadIDType

// IngressMsg the type for messages coming into the usecase/interactor layer
type IngressMsg entity.MsgBase

// EgressMsg the type of messages going out of the layer (has user generated metadata attached)
type EgressMsg entity.Msg

// MsgUsecase ifc
type MsgUsecase interface {
	// Checks whether the msg is valid before Enqueuing
	IsValid(msg *IngressMsg) (bool, error)

	// Enqueues the message into the system
	EnqueueMsg(msg *IngressMsg) (*EgressMsg, error)

	//For testing only should disallow an unqualified msg read in a real system, should be user based
	RetrieveMsg(mid MsgIDType) (*EgressMsg, error)
}

//NOTE MsgUsecase exposies the entity.Msg types at the usecase boundary.
//  It's not against Clean architecture since, the dependency is still inward, but
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
