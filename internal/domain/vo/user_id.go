package vo

import (
	"errors"
)

type UserID struct {
	value string
}

var (
	ErrUserIDInvalid = errors.New("invalid user ID")
)

func NewUserID(value string) UserID {
	return UserID{value: value}
}

func (id UserID) String() string {
	return id.value
}

func (id UserID) Equal(other UserID) bool {
	return id.value == other.value
}
