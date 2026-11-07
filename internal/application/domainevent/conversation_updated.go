package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationUpdated] = (*ConversationUpdatedHandler)(nil)

type ConversationUpdatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationUpdatedHandler(
	eventPublisher port.EventPublisher,
) *ConversationUpdatedHandler {
	return &ConversationUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationUpdatedHandler) Execute(ctx context.Context, domainEvent event.ConversationUpdated) error {
	integrationEvent := integrationevent.NewConversationUpdated(
		domainEvent.AggregateID(),
		domainEvent.UpdateTime(),
		domainEvent.Title().String(),
		domainEvent.Description().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
