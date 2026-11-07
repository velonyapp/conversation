package mysql

import (
	"context"
	"database/sql"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/infrastructure/messaging/protobuf"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ port.EventPublisher = (*eventPublisher)(nil)

type eventPublisher struct {
	db      *sql.DB
	encoder *protobuf.EventEncoder
}

func NewEventPublisher(
	db *sql.DB,
	encoder *protobuf.EventEncoder,
) port.EventPublisher {
	return &eventPublisher{
		db:      db,
		encoder: encoder,
	}
}

func (pub *eventPublisher) Publish(ctx context.Context, event integrationevent.IntegrationEvent) error {
	eventPayload, err := pub.encoder.Encode(event)
	if err != nil {
		return err
	}

	payload, err := protojson.Marshal(eventPayload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO outbox_events (
			id,
			type,
			aggregate_id,
			aggregate_type,
			payload,
			occur_time
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		event.ID(),
		event.Type(),
		event.AggregateID(),
		event.AggregateType(),
		string(payload),
		event.OccurTime(),
	)

	return err
}
