package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*MessageDeleted)(nil)

type MessageDeleted struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         *vo.UserID
	deleteTime     time.Time
}

func NewMessageDeleted(
	messageID vo.MessageID,
	conversationID vo.ConversationID,
	userID *vo.UserID,
	deleteTime time.Time,
) MessageDeleted {
	e := MessageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(messageID.String()),

		conversationID: conversationID,
		deleteTime:     deleteTime,
	}
	if userID != nil {
		value := *e.userID
		e.userID = &value
	}

	return e
}

func (e MessageDeleted) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e MessageDeleted) UserID() *vo.UserID {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}

func (e MessageDeleted) DeleteTime() time.Time {
	return e.deleteTime
}
