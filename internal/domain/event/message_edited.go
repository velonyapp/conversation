package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*MessageEdited)(nil)

type MessageEdited struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         *vo.UserID
	content        vo.MessageContent
	editTime       time.Time
	version        uint32
}

func NewMessageEdited(
	messageID vo.MessageID,
	conversationID vo.ConversationID,
	userID *vo.UserID,
	content vo.MessageContent,
	editTime time.Time,
	version uint32,
) MessageEdited {
	e := MessageEdited{
		BaseDomainEvent: NewBaseDomainEvent(messageID.String()),

		conversationID: conversationID,
		content:        content,
		editTime:       editTime,
		version:        version,
	}
	if userID != nil {
		value := *e.userID
		e.userID = &value
	}

	return e
}

func (e MessageEdited) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e MessageEdited) UserID() *vo.UserID {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}

func (e MessageEdited) Content() vo.MessageContent {
	return e.content
}

func (e MessageEdited) EditTime() time.Time {
	return e.editTime
}

func (e MessageEdited) Version() uint32 {
	return e.version
}
