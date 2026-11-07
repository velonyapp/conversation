package vo

import (
	"errors"

	"github.com/google/uuid"
)

type ConversationID struct {
	value string
}

var (
	ErrConversationIDInvalid = errors.New("invalid conversation ID")
)

func NewConversationID(value string) (ConversationID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return ConversationID{}, ErrConversationIDInvalid
	}

	return ConversationID{value: id.String()}, nil
}

func GenerateConversationID() ConversationID {
	return ConversationID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id ConversationID) String() string {
	return id.value
}

func (id ConversationID) Equal(other ConversationID) bool {
	return id.value == other.value
}
