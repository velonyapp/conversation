package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationCreated] = (*ConversationCreatedHandler)(nil)

type ConversationCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationCreatedHandler(
	eventPublisher port.EventPublisher,
) *ConversationCreatedHandler {
	return &ConversationCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationCreatedHandler) Execute(ctx context.Context, domainEvent event.ConversationCreated) error {
	integrationEvent := integrationevent.NewConversationCreated(
		domainEvent.AggregateID(),
		domainEvent.CreateTime(),
		domainEvent.Title().String(),
		domainEvent.Description().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
