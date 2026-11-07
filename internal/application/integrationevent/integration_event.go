package integrationevent

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationEvent interface {
	ID() string
	Type() string
	AggregateID() string
	AggregateType() string
	OccurTime() time.Time
}

type BaseIntegrationEvent struct {
	id          string
	aggregateID string
	occurTime   time.Time
}

func NewBaseIntegrationEvent(aggregateID string, occurTime time.Time) BaseIntegrationEvent {
	return BaseIntegrationEvent{
		id:          uuid.Must(uuid.NewV7()).String(),
		aggregateID: aggregateID,
		occurTime:   occurTime,
	}
}

func (e BaseIntegrationEvent) ID() string           { return e.id }
func (e BaseIntegrationEvent) AggregateID() string  { return e.aggregateID }
func (e BaseIntegrationEvent) OccurTime() time.Time { return e.occurTime }
