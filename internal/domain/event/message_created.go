package event

import (
	"time"

	"github.com/velonyapp/conversation/internal/domain/vo"
)

var _ DomainEvent = (*MessageCreated)(nil)

type MessageCreated struct {
	BaseDomainEvent

	conversationID vo.ConversationID
	userID         *vo.UserID
	content        vo.MessageContent
	createTime     time.Time
	version        uint32
}

func NewMessageCreated(
	messageID vo.MessageID,
	conversationID vo.ConversationID,
	userID *vo.UserID,
	content vo.MessageContent,
	createTime time.Time,
	version uint32,
) MessageCreated {
	e := MessageCreated{
		BaseDomainEvent: NewBaseDomainEvent(messageID.String()),

		conversationID: conversationID,
		content:        content,
		createTime:     createTime,
		version:        version,
	}
	if userID != nil {
		value := *e.userID
		e.userID = &value
	}

	return e
}

func (e MessageCreated) ConversationID() vo.ConversationID {
	return e.conversationID
}

func (e MessageCreated) UserID() *vo.UserID {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}

func (e MessageCreated) Content() vo.MessageContent {
	return e.content
}

func (e MessageCreated) CreateTime() time.Time {
	return e.createTime
}

func (e MessageCreated) Version() uint32 {
	return e.version
}
