package port

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
)

type EventPublisher interface {
	Publish(ctx context.Context, event integrationevent.IntegrationEvent) error
}
