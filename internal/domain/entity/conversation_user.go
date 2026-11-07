package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/conversation/internal/domain/event"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrConversationUserRemoved = errors.New("conversation user is removed")
)

type ConversationUser struct {
	id             vo.ConversationUserID
	conversationID vo.ConversationID
	userID         vo.UserID
	role           vo.ConversationUserRole
	createTime     time.Time
	updateTime     time.Time
	removeTime     *time.Time

	domainEvents []event.DomainEvent

	version uint32
	isNew   bool
	isDirty bool
}

func NewConversationUser(
	conversationID vo.ConversationID,
	userID vo.UserID,
	role vo.ConversationUserRole,
	now time.Time,
) *ConversationUser {
	c := &ConversationUser{
		id:             vo.GenerateConversationUserID(),
		conversationID: conversationID,
		userID:         userID,
		role:           role,
		createTime:     now,
		updateTime:     now,

		version: 0,
		isNew:   true,
		isDirty: true,
	}

	c.recordEvent(
		event.NewConversationUserJoined(
			c.id,
			c.conversationID,
			c.userID,
			c.role,
			c.createTime,
			c.version,
		),
	)

	return c
}

func ReconstituteConversationUser(
	id vo.ConversationUserID,
	conversationID vo.ConversationID,
	userID vo.UserID,
	role vo.ConversationUserRole,
	createTime time.Time,
	updateTime time.Time,
	removeTime *time.Time,
	version uint32,
) *ConversationUser {
	c := &ConversationUser{
		id:             id,
		conversationID: conversationID,
		userID:         userID,
		role:           role,
		createTime:     createTime,
		updateTime:     updateTime,
		version:        version,
	}
	if removeTime != nil {
		value := *removeTime
		c.removeTime = &value
	}

	return c
}

func (c *ConversationUser) ID() vo.ConversationUserID {
	return c.id
}

func (c *ConversationUser) ConversationID() vo.ConversationID {
	return c.conversationID
}

func (c *ConversationUser) UserID() vo.UserID {
	return c.userID
}

func (c *ConversationUser) Role() vo.ConversationUserRole {
	return c.role
}

func (c *ConversationUser) CreateTime() time.Time {
	return c.createTime
}

func (c *ConversationUser) UpdateTime() time.Time {
	return c.updateTime
}

func (c *ConversationUser) RemoveTime() *time.Time {
	if c.removeTime == nil {
		return nil
	}

	value := *c.removeTime
	return &value
}

func (c *ConversationUser) IsRemoved() bool {
	return c.removeTime != nil
}

func (c *ConversationUser) ChangeRole(newRole vo.ConversationUserRole, now time.Time) error {
	if c.IsRemoved() {
		return ErrConversationUserRemoved
	}
	if newRole.Equal(c.role) {
		return nil
	}

	c.role = newRole
	c.updateTime = now

	c.markDirty()

	c.recordEvent(
		event.NewConversationUserUpdated(
			c.id,
			c.conversationID,
			c.userID,
			c.role,
			c.updateTime,
			c.version,
		),
	)

	return nil
}

func (c *ConversationUser) Remove(now time.Time) error {
	if c.IsRemoved() {
		return nil
	}

	c.removeTime = &now
	c.updateTime = now

	c.markDirty()

	c.recordEvent(
		event.NewConversationUserRemoved(
			c.id,
			c.conversationID,
			c.userID,
			*c.removeTime,
		),
	)

	return nil
}

func (c *ConversationUser) PullEvents() []event.DomainEvent {
	pulled := c.domainEvents
	c.domainEvents = nil
	return pulled
}

func (c *ConversationUser) recordEvent(domainEvent event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, domainEvent)
}

func (c *ConversationUser) Version() uint32 {
	return c.version
}

func (c *ConversationUser) IsNew() bool {
	return c.isNew
}

func (c *ConversationUser) IsDirty() bool {
	return c.isDirty
}

func (c *ConversationUser) markDirty() {
	if c.isDirty {
		return
	}

	c.isDirty = true
	c.version++
}

func (c *ConversationUser) MarkPersisted() {
	c.isNew = false
	c.isDirty = false
}
