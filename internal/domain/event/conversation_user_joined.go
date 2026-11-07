package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationUserJoined)(nil)

type ConversationUserJoined struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         vo.UserID
	role           vo.ConversationUserRole
	joinTime       time.Time
	version        uint32
}

func NewConversationUserJoined(
	id vo.ConversationUserID,
	conversationID vo.ConversationID,
	userID vo.UserID,
	role vo.ConversationUserRole,
	joinTime time.Time,
	version uint32,
) ConversationUserJoined {
	return ConversationUserJoined{
		BaseDomainEvent: NewBaseDomainEvent(id.String()),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		joinTime:       joinTime,
		version:        version,
	}
}

func (e ConversationUserJoined) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e ConversationUserJoined) UserID() vo.UserID {
	return e.userID
}

func (e ConversationUserJoined) Role() vo.ConversationUserRole {
	return e.role
}

func (e ConversationUserJoined) JoinTime() time.Time {
	return e.joinTime
}

func (e ConversationUserJoined) Version() uint32 {
	return e.version
}
