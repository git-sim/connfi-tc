package entity

import (
	"time"
)

// MsgIDType set the id type
type MsgIDType uint64

// ThreadIDType set the thread id type
type ThreadIDType uint64

// MsgBase basic type coming into the system
type MsgBase struct {
	ParentMid   MsgIDType
	ScheduledAt time.Time
	SenderEmail string
	Recipients  []string
	Subject     string
	Body        []byte
}

// Msg type with system generated metadata attached more appropiate for storage
type Msg struct {
	Mid      MsgIDType
	Tid      ThreadIDType
	SentAt   time.Time
	SenderID AccountIDType
	M        MsgBase
}

// MsgEntry this is the decorated type used in Message Folders (inbox, archive, etc)
type MsgEntry struct {
	Mid       MsgIDType
	ViewedAt  time.Time
	IsRead    bool
	IsStarred bool
	M         Msg
}

// NewMsgEntry Creates a new MsgEntry from a Msg
func NewMsgEntry(msg Msg) *MsgEntry {
	return &MsgEntry{
		Mid: msg.Mid,
		M:   msg,
	}
}

// PendingMsgEntry for queued messages waiting for recipients. Early Optimization? could just scan the messages when a new user is added
type PendingMsgEntry struct {
	E              MsgEntry
	RecipientsLeft []string
}
