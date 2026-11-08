package domainevent

import (
	"context"
	"fmt"
	"reflect"

	"github.com/velonyapp/conversation/internal/domain/event"
)

type Dispatcher struct {
	handlers map[reflect.Type][]func(
		context.Context,
		event.DomainEvent,
	) error
}

func NewDispatcher(
	conversationCreatedHandler *ConversationCreatedHandler,
	conversationUpdatedHandler *ConversationUpdatedHandler,
	conversationDeletedHandler *ConversationDeletedHandler,
	conversationUserAddedHandler *ConversationUserAddedHandler,
	conversationUserUpdatedHandler *ConversationUserUpdatedHandler,
	conversationUserRemovedHandler *ConversationUserRemovedHandler,
	messageCreatedHandler *MessageCreatedHandler,
	messageEditedHandler *MessageEditedHandler,
	messageDeletedHandler *MessageDeletedHandler,
) *Dispatcher {
	dispatcher := &Dispatcher{
		handlers: make(
			map[reflect.Type][]func(
				context.Context,
				event.DomainEvent,
			) error,
		),
	}

	registerHandler(dispatcher, conversationCreatedHandler)
	registerHandler(dispatcher, conversationUpdatedHandler)
	registerHandler(dispatcher, conversationDeletedHandler)
	registerHandler(dispatcher, conversationUserAddedHandler)
	registerHandler(dispatcher, conversationUserUpdatedHandler)
	registerHandler(dispatcher, conversationUserRemovedHandler)
	registerHandler(dispatcher, messageCreatedHandler)
	registerHandler(dispatcher, messageEditedHandler)
	registerHandler(dispatcher, messageDeletedHandler)

	return dispatcher
}

func (d *Dispatcher) Dispatch(ctx context.Context, event event.DomainEvent) error {
	handlers := d.handlers[reflect.TypeOf(event)]

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func registerHandler[T event.DomainEvent](d *Dispatcher, h Handler[T]) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	d.handlers[eventType] = append(
		d.handlers[eventType],
		func(ctx context.Context, event event.DomainEvent) error {
			typedEvent, ok := event.(T)
			if !ok {
				return fmt.Errorf(
					"invalid domain event type: expected %v, got %T",
					eventType,
					event,
				)
			}

			return h.Execute(ctx, typedEvent)
		},
	)
}
