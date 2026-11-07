package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationDeleted] = (*ConversationDeletedHandler)(nil)

type ConversationDeletedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationDeletedHandler(
	eventPublisher port.EventPublisher,
) *ConversationDeletedHandler {
	return &ConversationDeletedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationDeletedHandler) Execute(ctx context.Context, domainEvent event.ConversationDeleted) error {
	integrationEvent := integrationevent.NewConversationDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
