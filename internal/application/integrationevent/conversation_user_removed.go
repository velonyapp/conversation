package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationUserRemoved)(nil)

type ConversationUserRemoved struct {
	BaseIntegrationEvent

	conversationID string
	userID         string
}

func NewConversationUserRemoved(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID string,
) ConversationUserRemoved {
	return ConversationUserRemoved{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
		userID:         userID,
	}
}

func (e ConversationUserRemoved) Type() string {
	return "conversation.conversation-user.removed"
}

func (e ConversationUserRemoved) AggregateType() string {
	return "conversation-user"
}

func (e ConversationUserRemoved) ConversationID() string {
	return e.conversationID
}

func (e ConversationUserRemoved) UserID() string {
	return e.userID
}
