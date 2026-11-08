package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationUserAdded)(nil)

type ConversationUserAdded struct {
	BaseIntegrationEvent

	conversationID string
	userID         string
	role           string
	version        uint32
}

func NewConversationUserAdded(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID string,
	role string,
	version uint32,
) ConversationUserAdded {
	return ConversationUserAdded{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		version:        version,
	}
}

func (e ConversationUserAdded) Type() string {
	return "conversation.conversation-user.added"
}

func (e ConversationUserAdded) AggregateType() string {
	return "conversation-user"
}

func (e ConversationUserAdded) ConversationID() string {
	return e.conversationID
}

func (e ConversationUserAdded) UserID() string {
	return e.userID
}

func (e ConversationUserAdded) Role() string {
	return e.role
}

func (e ConversationUserAdded) Version() uint32 {
	return e.version
}
