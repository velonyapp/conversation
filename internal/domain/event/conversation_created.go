package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*ConversationCreated)(nil)

type ConversationCreated struct {
	BaseDomainEvent

	title       vo.ConversationTitle
	description vo.ConversationDescription
	createTime  time.Time
	version     uint32
}

func NewConversationCreated(
	conversationID vo.ConversationID,
	title vo.ConversationTitle,
	description vo.ConversationDescription,
	createTime time.Time,
	version uint32,
) ConversationCreated {
	return ConversationCreated{
		BaseDomainEvent: NewBaseDomainEvent(conversationID.String()),

		title:       title,
		description: description,
		createTime:  createTime,
		version:     version,
	}
}

func (e ConversationCreated) Title() vo.ConversationTitle {
	return e.title
}

func (e ConversationCreated) Description() vo.ConversationDescription {
	return e.description
}

func (e ConversationCreated) CreateTime() time.Time {
	return e.createTime
}

func (e ConversationCreated) Version() uint32 {
	return e.version
}
