package common

import "time"

type ConversationResult struct {
	ID            string
	Title         string
	Description   string
	LastMessageID *string
	CreateTime    time.Time
	UpdateTime    time.Time
}

type ConversationUserResult struct {
	ID         string
	UserID     string
	Role       string
	CreateTime time.Time
	UpdateTime time.Time
}

type MessageResult struct {
	ID         string
	UserID     *string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}
