package mysql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/infrastructure/messaging/protobuf"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ port.EventPublisher = (*eventPublisher)(nil)

type eventPublisher struct {
	db      *sql.DB
	encoder *protobuf.EventEncoder
	log     *slog.Logger
}

func NewEventPublisher(
	db *sql.DB,
	encoder *protobuf.EventEncoder,
	log *slog.Logger,
) port.EventPublisher {
	return &eventPublisher{
		db:      db,
		encoder: encoder,
		log:     log,
	}
}

func (pub *eventPublisher) Publish(ctx context.Context, event integrationevent.IntegrationEvent) error {
	eventPayload, err := pub.encoder.Encode(event)
	if err != nil {
		pub.log.ErrorContext(ctx,
			"failed to encode integration event",
			slog.Any("event_id", event.ID()),
			slog.Any("event_type", event.Type()),
			slog.Any("error", err),
		)
		return err
	}

	payload, err := protojson.Marshal(eventPayload)
	if err != nil {
		pub.log.ErrorContext(ctx,
			"failed to marshal integration event payload",
			slog.Any("event_id", event.ID()),
			slog.Any("event_type", event.Type()),
			slog.Any("error", err),
		)
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
	if err != nil {
		pub.log.ErrorContext(ctx,
			"failed to insert integration event into outbox",
			slog.Any("event_id", event.ID()),
			slog.Any("event_type", event.Type()),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}
