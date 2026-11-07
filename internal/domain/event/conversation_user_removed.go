package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationUserRemoved)(nil)

type ConversationUserRemoved struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         vo.UserID
	removeTime     time.Time
}

func NewConversationUserRemoved(
	id vo.ConversationUserID,
	conversationID vo.ConversationID,
	userID vo.UserID,
	removeTime time.Time,
) ConversationUserRemoved {
	return ConversationUserRemoved{
		BaseDomainEvent: NewBaseDomainEvent(id.String()),

		conversationID: conversationID,
		userID:         userID,
		removeTime:     removeTime,
	}
}

func (e ConversationUserRemoved) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e ConversationUserRemoved) UserID() vo.UserID {
	return e.userID
}

func (e ConversationUserRemoved) RemoveTime() time.Time {
	return e.removeTime
}
