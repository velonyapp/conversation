package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.MessageEdited] = (*MessageEditedHandler)(nil)

type MessageEditedHandler struct {
	eventPublisher port.EventPublisher
}

func NewMessageEditedHandler(
	eventPublisher port.EventPublisher,
) *MessageEditedHandler {
	return &MessageEditedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *MessageEditedHandler) Execute(ctx context.Context, domainEvent event.MessageEdited) error {
	var userID *string
	if vo := domainEvent.UserID(); vo != nil {
		value := vo.String()
		userID = &value
	}

	integrationEvent := integrationevent.NewMessageEdited(
		domainEvent.AggregateID(),
		domainEvent.EditTime(),
		domainEvent.ConversationID().String(),
		userID,
		domainEvent.Content().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
