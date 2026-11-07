package vo

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrConversationDescriptionTooLong          = errors.New("conversation description must not exceed 2000 characters")
	ErrConversationDescriptionInvalidCharacter = errors.New("conversation description contains invalid characters")
)

type ConversationDescription struct {
	value string
}

func NewConversationDescription(value string) (ConversationDescription, error) {
	value = strings.TrimSpace(value)

	if !utf8.ValidString(value) {
		return ConversationDescription{}, ErrConversationDescriptionInvalidCharacter
	}

	if utf8.RuneCountInString(value) > 2000 {
		return ConversationDescription{}, ErrConversationDescriptionTooLong
	}

	for _, r := range value {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return ConversationDescription{}, ErrConversationDescriptionInvalidCharacter
		}
	}

	return ConversationDescription{value: value}, nil
}

func (u ConversationDescription) String() string {
	return u.value
}

func (u ConversationDescription) Equal(other ConversationDescription) bool {
	return u.value == other.value
}
