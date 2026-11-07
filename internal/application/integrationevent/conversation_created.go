package integrationevent

import "time"

var _ IntegrationEvent = (*ConversationCreated)(nil)

type ConversationCreated struct {
	BaseIntegrationEvent

	title       string
	description string
	version     uint32
}

func NewConversationCreated(
	aggregateID string,
	occurTime time.Time,
	title string,
	description string,
	version uint32,
) ConversationCreated {
	return ConversationCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		title:       title,
		description: description,
		version:     version,
	}
}

func (e ConversationCreated) Type() string {
	return "conversation.conversation.created"
}

func (e ConversationCreated) AggregateType() string {
	return "conversation"
}

func (e ConversationCreated) Title() string {
	return e.title
}

func (e ConversationCreated) Description() string {
	return e.description
}

func (e ConversationCreated) Version() uint32 {
	return e.version
}
