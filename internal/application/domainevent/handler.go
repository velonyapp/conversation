package domainevent

import (
	"context"

	"github.com/velonyapp/conversation/internal/domain/event"
)

type Handler[T event.DomainEvent] interface {
	Execute(ctx context.Context, event T) error
}
