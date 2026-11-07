package integrationevent

import "time"

var _ IntegrationEvent = (*MessageEdited)(nil)

type MessageEdited struct {
	BaseIntegrationEvent

	conversationID string
	userID         *string
	content        string
	version        uint32
}

func NewMessageEdited(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID *string,
	content string,
	version uint32,
) MessageEdited {
	e := MessageEdited{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
		content:        content,
		version:        version,
	}
	if userID != nil {
		value := *userID
		e.userID = &value
	}

	return e
}

func (e MessageEdited) Type() string {
	return "conversation.message.edited"
}

func (e MessageEdited) AggregateType() string {
	return "message"
}

func (e MessageEdited) ConversationID() string {
	return e.conversationID
}

func (e MessageEdited) UserID() *string {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}

func (e MessageEdited) Content() string {
	return e.content
}

func (e MessageEdited) Version() uint32 {
	return e.version
}
