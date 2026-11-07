package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.MessageCreated] = (*MessageCreatedHandler)(nil)

type MessageCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewMessageCreatedHandler(
	eventPublisher port.EventPublisher,
) *MessageCreatedHandler {
	return &MessageCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *MessageCreatedHandler) Execute(ctx context.Context, domainEvent event.MessageCreated) error {
	var userID *string
	if vo := domainEvent.UserID(); vo != nil {
		value := vo.String()
		userID = &value
	}

	integrationEvent := integrationevent.NewMessageCreated(
		domainEvent.AggregateID(),
		domainEvent.CreateTime(),
		domainEvent.ConversationID().String(),
		userID,
		domainEvent.Content().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
