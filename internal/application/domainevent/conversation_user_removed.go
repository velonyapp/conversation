package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationUserRemoved] = (*ConversationUserRemovedHandler)(nil)

type ConversationUserRemovedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationUserRemovedHandler(
	eventPublisher port.EventPublisher,
) *ConversationUserRemovedHandler {
	return &ConversationUserRemovedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationUserRemovedHandler) Execute(ctx context.Context, domainEvent event.ConversationUserRemoved) error {
	integrationEvent := integrationevent.NewConversationUserRemoved(
		domainEvent.AggregateID(),
		domainEvent.RemoveTime(),
		domainEvent.ConversationID().String(),
		domainEvent.UserID().String(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
