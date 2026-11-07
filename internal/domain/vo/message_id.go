package vo

import (
	"errors"

	"github.com/google/uuid"
)

type MessageID struct {
	value string
}

var (
	ErrMessageIDInvalid = errors.New("invalid message ID")
)

func NewMessageID(value string) (MessageID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return MessageID{}, ErrMessageIDInvalid
	}

	return MessageID{value: id.String()}, nil
}

func GenerateMessageID() MessageID {
	return MessageID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id MessageID) String() string {
	return id.value
}

func (id MessageID) Equal(other MessageID) bool {
	return id.value == other.value
}
