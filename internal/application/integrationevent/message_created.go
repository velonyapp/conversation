package integrationevent

import "time"

var _ IntegrationEvent = (*MessageCreated)(nil)

type MessageCreated struct {
	BaseIntegrationEvent

	conversationID string
	userID         *string
	content        string
	version        uint32
}

func NewMessageCreated(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID *string,
	content string,
	version uint32,
) MessageCreated {
	e := MessageCreated{
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

func (e MessageCreated) Type() string {
	return "conversation.message.created"
}

func (e MessageCreated) AggregateType() string {
	return "message"
}

func (e MessageCreated) ConversationID() string {
	return e.conversationID
}

func (e MessageCreated) UserID() *string {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}

func (e MessageCreated) Content() string {
	return e.content
}

func (e MessageCreated) Version() uint32 {
	return e.version
}
