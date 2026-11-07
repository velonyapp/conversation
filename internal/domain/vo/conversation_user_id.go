package vo

import (
	"errors"

	"github.com/google/uuid"
)

type ConversationUserID struct {
	value string
}

var (
	ErrConversationUserIDInvalid = errors.New("invalid conversation user ID")
)

func NewConversationUserID(value string) (ConversationUserID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return ConversationUserID{}, ErrConversationUserIDInvalid
	}

	return ConversationUserID{value: id.String()}, nil
}

func GenerateConversationUserID() ConversationUserID {
	return ConversationUserID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id ConversationUserID) String() string {
	return id.value
}

func (id ConversationUserID) Equal(other ConversationUserID) bool {
	return id.value == other.value
}
