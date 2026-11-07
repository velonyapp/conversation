package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationUserJoined)(nil)

type ConversationUserJoined struct {
	BaseIntegrationEvent

	conversationID string
	userID         string
	role           string
	version        uint32
}

func NewConversationUserJoined(
	aggregateID string,
	occurTime time.Time,
	conversationID string,
	userID string,
	role string,
	version uint32,
) ConversationUserJoined {
	return ConversationUserJoined{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		conversationID: conversationID,
		userID:         userID,
		role:           role,
		version:        version,
	}
}

func (e ConversationUserJoined) Type() string {
	return "conversation.conversation-user.joined"
}

func (e ConversationUserJoined) AggregateType() string {
	return "conversation-user"
}

func (e ConversationUserJoined) ConversationID() string {
	return e.conversationID
}

func (e ConversationUserJoined) UserID() string {
	return e.userID
}

func (e ConversationUserJoined) Role() string {
	return e.role
}

func (e ConversationUserJoined) Version() uint32 {
	return e.version
}
