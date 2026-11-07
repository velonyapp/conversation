package vo

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	ErrMessageContentTooLong          = errors.New("message content must not exceed 65536 characters")
	ErrMessageContentInvalidCharacter = errors.New("message content contains invalid characters")
)

type MessageContent struct {
	value string
}

func NewMessageContent(value string) (MessageContent, error) {
	if !utf8.ValidString(value) {
		return MessageContent{}, ErrMessageContentInvalidCharacter
	}

	if utf8.RuneCountInString(value) > 65536 {
		return MessageContent{}, ErrMessageContentTooLong
	}

	for _, r := range value {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return MessageContent{}, ErrMessageContentInvalidCharacter
		}
	}

	return MessageContent{value: value}, nil
}

func (u MessageContent) String() string {
	return u.value
}

func (u MessageContent) Equal(other MessageContent) bool {
	return u.value == other.value
}
