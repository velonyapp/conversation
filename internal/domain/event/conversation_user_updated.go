package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationUserUpdated)(nil)

type ConversationUserUpdated struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         vo.UserID
	role           vo.ConversationUserRole
	updateTime     time.Time
	version        uint32
}

func NewConversationUserUpdated(
	id vo.ConversationUserID,
	conversationID vo.ConversationID,
	userID vo.UserID,
	role vo.ConversationUserRole,
	updateTime time.Time,
	version uint32,
) ConversationUserUpdated {
	return ConversationUserUpdated{
		BaseDomainEvent: NewBaseDomainEvent(id.String()),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		updateTime:     updateTime,
		version:        version,
	}
}

func (e ConversationUserUpdated) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e ConversationUserUpdated) UserID() vo.UserID {
	return e.userID
}

func (e ConversationUserUpdated) Role() vo.ConversationUserRole {
	return e.role
}

func (e ConversationUserUpdated) UpdateTime() time.Time {
	return e.updateTime
}

func (e ConversationUserUpdated) Version() uint32 {
	return e.version
}
