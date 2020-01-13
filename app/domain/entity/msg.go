package entity

import (
	"time"
)

// MsgIDType set the id type
type MsgIDType uint64

const MsgIDBits = 64
const MsgIDStringBase = 16

// ThreadIDType set the thread id type
type ThreadIDType uint64

const ThreadIDBits = 64
const ThreadIDStringBase = 16

// MsgBase basic type coming into the system
type MsgBase struct {
	ParentMid   MsgIDType `json:"parentmsgid"`
	CreatedAt   time.Time `json:"createdat"`
	ScheduledAt time.Time `json:"scheduledat"`
	SenderEmail string    `json:"senderemail"`
	Recipients  []string  `json:"recipients"`
	Subject     string    `json:"subject"`
	Body        []byte    `json:"body"`
}

// Msg type with system generated metadata attached more appropiate for storage
type Msg struct {
	Mid      MsgIDType     `json:"msgid"`
	Tid      ThreadIDType  `json:"threadid"`
	SentAt   time.Time     `json:"sentat"`
	SenderID AccountIDType `json:"senderid"`
	M        MsgBase
}

// MsgEntry this is the decorated type used in Message Folders (inbox, archive, etc)
type MsgEntry struct {
	Mid       MsgIDType `json:"msgid"`
	ViewedAt  time.Time `json:"viewedat"`
	IsViewed  bool      `json:"isviewed"`
	IsStarred bool      `json:"isstarred"`
	Folder    string    `json:"foldername"`
	M         Msg
}

func NewMsg(msgbase MsgBase) *Msg {
	return &Msg{
		M: msgbase,
	}
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
	E             MsgEntry
	RecipientLeft string
}

// NewPendingMsgEntry Creates a new PendingMsgEntry from a MsgEntry
func NewPendingMsgEntry(me MsgEntry, remail string) *PendingMsgEntry {
	return &PendingMsgEntry{
		E:             me,
		RecipientLeft: remail,
	}
}
