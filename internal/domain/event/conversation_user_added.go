package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationUserAdded)(nil)

type ConversationUserAdded struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         vo.UserID
	role           vo.ConversationUserRole
	addTime        time.Time
	version        uint32
}

func NewConversationUserAdded(
	id vo.ConversationUserID,
	conversationID vo.ConversationID,
	userID vo.UserID,
	role vo.ConversationUserRole,
	addTime time.Time,
	version uint32,
) ConversationUserAdded {
	return ConversationUserAdded{
		BaseDomainEvent: NewBaseDomainEvent(id.String()),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		addTime:        addTime,
		version:        version,
	}
}

func (e ConversationUserAdded) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e ConversationUserAdded) UserID() vo.UserID {
	return e.userID
}

func (e ConversationUserAdded) Role() vo.ConversationUserRole {
	return e.role
}

func (e ConversationUserAdded) AddTime() time.Time {
	return e.addTime
}

func (e ConversationUserAdded) Version() uint32 {
	return e.version
}
