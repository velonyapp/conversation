package protobuf

import (
	"fmt"

	"github.com/velonyapp/conversation/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
)

type EventEncoder struct{}

func NewEventEncoder() *EventEncoder {
	return &EventEncoder{}
}

func (e *EventEncoder) Encode(eventPayload integrationevent.IntegrationEvent) (proto.Message, error) {
	switch eventPayload := eventPayload.(type) {
	case integrationevent.ConversationCreated:
		return conversationCreatedPayload(eventPayload), nil
	case integrationevent.ConversationUpdated:
		return conversationUpdatedPayload(eventPayload), nil
	case integrationevent.ConversationDeleted:
		return conversationDeletedPayload(eventPayload), nil
	case integrationevent.ConversationUserJoined:
		return conversationUserJoinedPayload(eventPayload), nil
	case integrationevent.ConversationUserUpdated:
		return conversationUserUpdatedPayload(eventPayload), nil
	case integrationevent.ConversationUserRemoved:
		return conversationUserRemovedPayload(eventPayload), nil
	case integrationevent.MessageCreated:
		return messageCreatedPayload(eventPayload), nil
	case integrationevent.MessageEdited:
		return messageEditedPayload(eventPayload), nil
	case integrationevent.MessageDeleted:
		return messageDeletedPayload(eventPayload), nil
	default:
		return nil, fmt.Errorf("unknown integration event %T", eventPayload)
	}
}
