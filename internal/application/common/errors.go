package common

import "errors"

var (
	ErrConversationNotFound     = errors.New("conversation not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrConversationUserNotFound = errors.New("conversation user not found")
	ErrMessageNotFound          = errors.New("message not found")
)
