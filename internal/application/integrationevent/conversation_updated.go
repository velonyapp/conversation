package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationUpdated)(nil)

type ConversationUpdated struct {
	BaseIntegrationEvent

	title       string
	description string
	version     uint32
}

func NewConversationUpdated(
	aggregateID string,
	occurTime time.Time,
	title string,
	description string,
	version uint32,
) ConversationUpdated {
	return ConversationUpdated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		title:       title,
		description: description,
		version:     version,
	}
}

func (e ConversationUpdated) Type() string {
	return "conversation.conversation.updated"
}

func (e ConversationUpdated) AggregateType() string {
	return "conversation"
}

func (e ConversationUpdated) Title() string {
	return e.title
}

func (e ConversationUpdated) Description() string {
	return e.description
}

func (e ConversationUpdated) Version() uint32 {
	return e.version
}
