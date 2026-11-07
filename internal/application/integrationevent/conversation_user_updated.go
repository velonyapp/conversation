package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationUserUpdated)(nil)

type ConversationUserUpdated struct {
	BaseIntegrationEvent

	conversationID string
	userID         string
	role           string
	version        uint32
}

func NewConversationUserUpdated(
	aggregateID string,
	occurTime time.Time,

	conversationID string,
	userID string,
	role string,
	version uint32,
) ConversationUserUpdated {
	return ConversationUserUpdated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		version:        version,
	}
}

func (e ConversationUserUpdated) Type() string {
	return "conversation.conversation-user.updated"
}

func (e ConversationUserUpdated) AggregateType() string {
	return "conversation-user"
}

func (e ConversationUserUpdated) ConversationID() string {
	return e.conversationID
}

func (e ConversationUserUpdated) UserID() string {
	return e.userID
}

func (e ConversationUserUpdated) Role() string {
	return e.role
}

func (e ConversationUserUpdated) Version() uint32 {
	return e.version
}
