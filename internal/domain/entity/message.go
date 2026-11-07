package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/conversation/internal/domain/event"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrMessageDeleted = errors.New("message is deleted")
)

type Message struct {
	id             vo.MessageID
	conversationID vo.ConversationID
	userID         *vo.UserID
	content        vo.MessageContent
	createTime     time.Time
	updateTime     time.Time
	deleteTime     *time.Time

	domainEvents []event.DomainEvent

	version uint32
	isNew   bool
	isDirty bool
}

func NewMessage(
	conversationID vo.ConversationID,
	userID *vo.UserID,
	content vo.MessageContent,
	now time.Time,
) *Message {
	m := &Message{
		id:             vo.GenerateMessageID(),
		conversationID: conversationID,
		content:        content,
		createTime:     now,
		updateTime:     now,

		version: 0,
		isNew:   true,
		isDirty: true,
	}
	if userID != nil {
		value := *userID
		m.userID = &value
	}

	m.recordEvent(
		event.NewMessageCreated(
			m.id,
			m.conversationID,
			m.userID,
			m.content,
			m.createTime,
			m.version,
		),
	)

	return m
}

func ReconstituteMessage(
	id vo.MessageID,
	conversationID vo.ConversationID,
	userID *vo.UserID,
	content vo.MessageContent,
	createTime time.Time,
	updateTime time.Time,
	deleteTime *time.Time,
	version uint32,
) *Message {
	m := &Message{
		id:             id,
		conversationID: conversationID,
		content:        content,
		createTime:     createTime,
		updateTime:     updateTime,
		version:        version,
	}
	if userID != nil {
		value := *userID
		m.userID = &value
	}
	if deleteTime != nil {
		value := *deleteTime
		m.deleteTime = &value
	}

	return m
}

func (m *Message) ID() vo.MessageID {
	return m.id
}

func (m *Message) ConversationID() vo.ConversationID {
	return m.conversationID
}

func (m *Message) UserID() *vo.UserID {
	if m.userID == nil {
		return nil
	}

	value := *m.userID
	return &value
}

func (m *Message) Content() vo.MessageContent {
	return m.content
}

func (m *Message) CreateTime() time.Time {
	return m.createTime
}

func (m *Message) UpdateTime() time.Time {
	return m.updateTime
}

func (m *Message) DeleteTime() *time.Time {
	if m.deleteTime == nil {
		return nil
	}

	value := *m.deleteTime
	return &value
}

func (m *Message) HasUser() bool {
	return m.userID != nil
}

func (m *Message) IsDeleted() bool {
	return m.deleteTime != nil
}

func (m *Message) Edit(newContent vo.MessageContent, now time.Time) error {
	if m.IsDeleted() {
		return ErrMessageDeleted
	}
	if newContent.Equal(m.content) {
		return nil
	}

	m.content = newContent
	m.updateTime = now

	m.markDirty()

	m.recordEvent(
		event.NewMessageEdited(
			m.id,
			m.conversationID,
			m.userID,
			m.content,
			m.updateTime,
			m.version,
		),
	)

	return nil
}

func (m *Message) Delete(now time.Time) error {
	if m.IsDeleted() {
		return nil
	}

	m.updateTime = now
	m.deleteTime = &now

	m.markDirty()

	m.recordEvent(
		event.NewMessageDeleted(
			m.id,
			m.conversationID,
			m.userID,
			*m.deleteTime,
		),
	)

	return nil
}

func (m *Message) PullEvents() []event.DomainEvent {
	pulled := m.domainEvents
	m.domainEvents = nil
	return pulled
}

func (m *Message) recordEvent(domainEvent event.DomainEvent) {
	m.domainEvents = append(m.domainEvents, domainEvent)
}

func (m *Message) Version() uint32 {
	return m.version
}

func (m *Message) IsNew() bool {
	return m.isNew
}

func (m *Message) IsDirty() bool {
	return m.isDirty
}

func (m *Message) markDirty() {
	if m.isDirty {
		return
	}

	m.isDirty = true
	m.version++
}

func (m *Message) MarkPersisted() {
	m.isNew = false
	m.isDirty = false
}
