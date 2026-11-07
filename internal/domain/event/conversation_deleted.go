package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationDeleted)(nil)

type ConversationDeleted struct {
	BaseDomainEvent

	deleteTime time.Time
}

func NewConversationDeleted(
	conversationID vo.ConversationID,
	deleteTime time.Time,
	version uint32,
) ConversationDeleted {
	return ConversationDeleted{
		BaseDomainEvent: NewBaseDomainEvent(conversationID.String()),
		deleteTime:      deleteTime,
	}
}

func (e ConversationDeleted) DeleteTime() time.Time {
	return e.deleteTime
}
