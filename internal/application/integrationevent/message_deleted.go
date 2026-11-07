package integrationevent

import "time"

var _ IntegrationEvent = (*MessageDeleted)(nil)

type MessageDeleted struct {
	BaseIntegrationEvent

	conversationID string
	userID         *string
}

func NewMessageDeleted(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID *string,
) MessageDeleted {
	e := MessageDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
	}
	if userID != nil {
		value := *userID
		e.userID = &value
	}

	return e
}

func (e MessageDeleted) Type() string {
	return "conversation.message.deleted"
}

func (e MessageDeleted) AggregateType() string {
	return "message"
}

func (e MessageDeleted) ConversationID() string {
	return e.conversationID
}

func (e MessageDeleted) UserID() *string {
	if e.userID == nil {
		return nil
	}

	value := *e.userID
	return &value
}
