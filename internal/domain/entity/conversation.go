package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/conversation/internal/domain/event"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrConversationDeleted = errors.New("conversation is deleted")
)

type Conversation struct {
	id          vo.ConversationID
	title       vo.ConversationTitle
	description vo.ConversationDescription
	createTime  time.Time
	updateTime  time.Time
	deleteTime  *time.Time

	domainEvents []event.DomainEvent

	version uint32
	isNew   bool
	isDirty bool
}

func NewConversation(
	title vo.ConversationTitle,
	description vo.ConversationDescription,
	now time.Time,
) *Conversation {
	c := &Conversation{
		id:          vo.GenerateConversationID(),
		title:       title,
		description: description,
		createTime:  now,
		updateTime:  now,

		version: 0,
		isNew:   true,
		isDirty: true,
	}

	c.recordEvent(
		event.NewConversationCreated(
			c.id,
			c.title,
			c.description,
			c.createTime,
			c.version,
		),
	)

	return c
}

func ReconstituteConversation(
	id vo.ConversationID,
	title vo.ConversationTitle,
	description vo.ConversationDescription,
	createTime time.Time,
	updateTime time.Time,
	deleteTime *time.Time,
	version uint32,
) *Conversation {
	c := &Conversation{
		id:          id,
		title:       title,
		description: description,
		createTime:  createTime,
		updateTime:  updateTime,
		version:     version,
	}
	if deleteTime != nil {
		value := *deleteTime
		c.deleteTime = &value
	}

	return c
}

func (c *Conversation) ID() vo.ConversationID {
	return c.id
}

func (c *Conversation) Title() vo.ConversationTitle {
	return c.title
}

func (c *Conversation) Description() vo.ConversationDescription {
	return c.description
}

func (c *Conversation) CreateTime() time.Time {
	return c.createTime
}

func (c *Conversation) UpdateTime() time.Time {
	return c.updateTime
}

func (c *Conversation) DeleteTime() *time.Time {
	if c.deleteTime == nil {
		return nil
	}

	value := *c.deleteTime
	return &value
}

func (c *Conversation) IsDeleted() bool {
	return c.deleteTime != nil
}

func (c *Conversation) ChangeTitle(newTitle vo.ConversationTitle, now time.Time) error {
	if c.IsDeleted() {
		return ErrConversationDeleted
	}
	if newTitle.Equal(c.title) {
		return nil
	}

	c.title = newTitle
	c.updateTime = now

	c.markDirty()

	c.recordEvent(
		event.NewConversationUpdated(
			c.id,
			c.title,
			c.description,
			c.updateTime,
			c.version,
		),
	)

	return nil
}

func (c *Conversation) ChangeDescription(newDescription vo.ConversationDescription, now time.Time) error {
	if c.IsDeleted() {
		return ErrConversationDeleted
	}
	if newDescription.Equal(c.description) {
		return nil
	}

	c.description = newDescription
	c.updateTime = now

	c.markDirty()

	c.recordEvent(
		event.NewConversationUpdated(
			c.id,
			c.title,
			c.description,
			c.updateTime,
			c.version,
		),
	)

	return nil
}

func (c *Conversation) Delete(now time.Time) error {
	if c.IsDeleted() {
		return nil
	}

	c.deleteTime = &now

	c.markDirty()

	c.recordEvent(
		event.NewConversationDeleted(
			c.id,
			*c.deleteTime,
			c.version,
		),
	)

	return nil
}

func (c *Conversation) PullEvents() []event.DomainEvent {
	pulled := c.domainEvents
	c.domainEvents = nil
	return pulled
}

func (c *Conversation) recordEvent(domainEvent event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, domainEvent)
}

func (c *Conversation) Version() uint32 {
	return c.version
}

func (c *Conversation) IsNew() bool {
	return c.isNew
}

func (c *Conversation) IsDirty() bool {
	return c.isDirty
}

func (c *Conversation) markDirty() {
	if c.isDirty {
		return
	}

	c.isDirty = true
	c.version++
}

func (c *Conversation) MarkPersisted() {
	c.isNew = false
	c.isDirty = false
}
