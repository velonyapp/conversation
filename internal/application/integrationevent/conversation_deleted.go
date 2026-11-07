package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationDeleted)(nil)

type ConversationDeleted struct {
	BaseIntegrationEvent
}

func NewConversationDeleted(
	aggregateID string,
	occurTime time.Time,
) ConversationDeleted {
	return ConversationDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),
	}
}

func (e ConversationDeleted) Type() string {
	return "conversation.conversation.deleted"
}

func (e ConversationDeleted) AggregateType() string {
	return "conversation"
}
