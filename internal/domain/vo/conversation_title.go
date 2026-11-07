package vo

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrConversationTitleEmpty            = errors.New("conversation title must not be empty")
	ErrConversationTitleTooLong          = errors.New("conversation title must not exceed 100 characters")
	ErrConversationTitleInvalidCharacter = errors.New("conversation title contains invalid characters")
)

type ConversationTitle struct {
	value string
}

func NewConversationTitle(value string) (ConversationTitle, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return ConversationTitle{}, ErrConversationTitleEmpty
	}

	if !utf8.ValidString(value) {
		return ConversationTitle{}, ErrConversationTitleInvalidCharacter
	}

	if utf8.RuneCountInString(value) > 100 {
		return ConversationTitle{}, ErrConversationTitleTooLong
	}

	for _, r := range value {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			return ConversationTitle{}, ErrConversationTitleInvalidCharacter
		}
	}

	return ConversationTitle{value: value}, nil
}

func (u ConversationTitle) String() string {
	return u.value
}

func (u ConversationTitle) Equal(other ConversationTitle) bool {
	return u.value == other.value
}
