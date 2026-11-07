package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationUpdated)(nil)

type ConversationUpdated struct {
	BaseDomainEvent

	title       vo.ConversationTitle
	description vo.ConversationDescription
	updateTime  time.Time
	version     uint32
}

func NewConversationUpdated(
	conversationID vo.ConversationID,
	title vo.ConversationTitle,
	description vo.ConversationDescription,
	updateTime time.Time,
	version uint32,
) ConversationUpdated {
	return ConversationUpdated{
		BaseDomainEvent: NewBaseDomainEvent(conversationID.String()),

		title:       title,
		description: description,
		updateTime:  updateTime,
		version:     version,
	}
}

func (e ConversationUpdated) Title() vo.ConversationTitle {
	return e.title
}

func (e ConversationUpdated) Description() vo.ConversationDescription {
	return e.description
}

func (e ConversationUpdated) UpdateTime() time.Time {
	return e.updateTime
}

func (e ConversationUpdated) Version() uint32 {
	return e.version
}
