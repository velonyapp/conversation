package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/event"
)

var _ Handler[event.MessageDeleted] = (*MessageDeletedHandler)(nil)

type MessageDeletedHandler struct {
	eventPublisher port.EventPublisher
}

func NewMessageDeletedHandler(
	eventPublisher port.EventPublisher,
) *MessageDeletedHandler {
	return &MessageDeletedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *MessageDeletedHandler) Execute(ctx context.Context, domainEvent event.MessageDeleted) error {
	var userID *string
	if vo := domainEvent.UserID(); vo != nil {
		value := vo.String()
		userID = &value
	}

	integrationEvent := integrationevent.NewMessageDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime(),
		domainEvent.ConversationID().String(),
		userID,
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
