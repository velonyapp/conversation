package event

import "github.com/google/uuid"

type DomainEvent interface {
	ID() string
	AggregateID() string
}

type BaseDomainEvent struct {
	id          string
	aggregateID string
}

func NewBaseDomainEvent(aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		id:          uuid.Must(uuid.NewV7()).String(),
		aggregateID: aggregateID,
	}
}

func (e BaseDomainEvent) ID() string          { return e.id }
func (e BaseDomainEvent) AggregateID() string { return e.aggregateID }
