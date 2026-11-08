package service

import (
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
)

var (
	ErrCannotActOnSelf          = errors.New("cannot perform this action on yourself")
	ErrConversationMismatch     = errors.New("conversation mismatch")
	ErrCannotElectOther         = errors.New("cannot elect another user")
	ErrCannotPromoteOther       = errors.New("cannot promote another user")
	ErrCannotDemoteOther        = errors.New("cannot demote another user")
	ErrCannotKickOther          = errors.New("cannot kick another user")
	ErrCannotUpdateConversation = errors.New("cannot update conversation")
	ErrCannotDeleteConversation = errors.New("cannot delete conversation")
	ErrCannotLeaveConversation  = errors.New("cannot leave conversation")
	ErrCannotEditMessage        = errors.New("cannot edit message")
	ErrCannotDeleteMessage      = errors.New("cannot delete message")
)

type ConversationUserPolicy struct {
}

func NewConversationUserPolicy() *ConversationUserPolicy {
	return &ConversationUserPolicy{}
}

func (p *ConversationUserPolicy) CanElectOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return ErrCannotActOnSelf
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return ErrConversationMismatch
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}

	return ErrCannotElectOther
}

func (p *ConversationUserPolicy) CanPromoteOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return ErrCannotActOnSelf
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return ErrConversationMismatch
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return ErrCannotPromoteOther
}

func (p *ConversationUserPolicy) CanDemoteOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return ErrCannotActOnSelf
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return ErrConversationMismatch
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return ErrCannotDemoteOther
}

func (p *ConversationUserPolicy) CanKickOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return ErrCannotActOnSelf
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return ErrConversationMismatch
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return ErrCannotKickOther
}

func (p *ConversationUserPolicy) CanAddUser(conversationUser *entity.ConversationUser) error {
	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() {
		return nil
	}

	return ErrCannotKickOther
}

func (p *ConversationUserPolicy) CanUpdateConversation(conversationUser *entity.ConversationUser) error {
	if conversationUser.Role().IsOwner() || conversationUser.Role().IsAdmin() {
		return nil
	}

	return ErrCannotUpdateConversation
}

func (p *ConversationUserPolicy) CanDeleteConversation(conversationUser *entity.ConversationUser) error {
	if conversationUser.Role().IsOwner() {
		return nil
	}

	return ErrCannotDeleteConversation
}

func (p *ConversationUserPolicy) CanLeaveConversation(conversationUser *entity.ConversationUser) error {
	if conversationUser.Role().IsOwner() {
		return ErrCannotLeaveConversation
	}

	return nil
}

func (p *ConversationUserPolicy) CanSendMessage(conversationUser *entity.ConversationUser) error {
	return nil
}

func (p *ConversationUserPolicy) CanEditMessage(conversationUser *entity.ConversationUser, message *entity.Message) error {
	if !conversationUser.ConversationID().Equal(message.ConversationID()) {
		return ErrConversationMismatch
	}

	if message.HasUser() && conversationUser.UserID().Equal(*message.UserID()) {
		return nil
	}

	return ErrCannotEditMessage
}

func (p *ConversationUserPolicy) CanDeleteMessage(conversationUser *entity.ConversationUser, message *entity.Message) error {
	if !conversationUser.ConversationID().Equal(message.ConversationID()) {
		return ErrConversationMismatch
	}

	if message.HasUser() && conversationUser.UserID().Equal(*message.UserID()) {
		return nil
	}
	if conversationUser.Role().IsOwner() || conversationUser.Role().IsAdmin() {
		return nil
	}

	return ErrCannotDeleteMessage
}
