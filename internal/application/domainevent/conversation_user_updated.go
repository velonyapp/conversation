package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationUserUpdated] = (*ConversationUserUpdatedHandler)(nil)

type ConversationUserUpdatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationUserUpdatedHandler(
	eventPublisher port.EventPublisher,
) *ConversationUserUpdatedHandler {
	return &ConversationUserUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationUserUpdatedHandler) Execute(ctx context.Context, domainEvent event.ConversationUserUpdated) error {
	integrationEvent := integrationevent.NewConversationUserUpdated(
		domainEvent.AggregateID(),
		domainEvent.UpdateTime(),
		domainEvent.ConversationID().String(),
		domainEvent.UserID().String(),
		domainEvent.Role().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
