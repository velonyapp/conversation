package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationUserAdded] = (*ConversationUserAddedHandler)(nil)

type ConversationUserAddedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationUserAddedHandler(
	eventPublisher port.EventPublisher,
) *ConversationUserAddedHandler {
	return &ConversationUserAddedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationUserAddedHandler) Execute(ctx context.Context, domainEvent event.ConversationUserAdded) error {
	integrationEvent := integrationevent.NewConversationUserAdded(
		domainEvent.AggregateID(),
		domainEvent.AddTime(),
		domainEvent.ConversationID().String(),
		domainEvent.UserID().String(),
		domainEvent.Role().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
