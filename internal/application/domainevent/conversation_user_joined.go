package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.ConversationUserJoined] = (*ConversationUserJoinedHandler)(nil)

type ConversationUserJoinedHandler struct {
	eventPublisher port.EventPublisher
}

func NewConversationUserJoinedHandler(
	eventPublisher port.EventPublisher,
) *ConversationUserJoinedHandler {
	return &ConversationUserJoinedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ConversationUserJoinedHandler) Execute(ctx context.Context, domainEvent event.ConversationUserJoined) error {
	integrationEvent := integrationevent.NewConversationUserJoined(
		domainEvent.AggregateID(),
		domainEvent.JoinTime(),
		domainEvent.ConversationID().String(),
		domainEvent.UserID().String(),
		domainEvent.Role().String(),
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
