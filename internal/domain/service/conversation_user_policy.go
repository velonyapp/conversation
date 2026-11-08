package service

import (
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
)

var ()

type ConversationUserPolicy struct {
}

func NewConversationUserPolicy() *ConversationUserPolicy {
	return &ConversationUserPolicy{}
}

func (p *ConversationUserPolicy) CanElectOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return errors.New("")
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanPromoteOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return errors.New("")
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanDemoteOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return errors.New("")
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanKickOther(conversationUser, other *entity.ConversationUser) error {
	if conversationUser.ID().Equal(other.ID()) {
		return errors.New("")
	}
	if !conversationUser.ConversationID().Equal(other.ConversationID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}
	if conversationUser.Role().IsAdmin() && !other.Role().IsOwner() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanUpdateConversation(conversationUser *entity.ConversationUser, conversation *entity.Conversation) error {
	if !conversationUser.ConversationID().Equal(conversation.ID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() || conversationUser.Role().IsAdmin() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanDeleteConversation(conversationUser *entity.ConversationUser, conversation *entity.Conversation) error {
	if !conversationUser.ConversationID().Equal(conversation.ID()) {
		return errors.New("")
	}

	if conversationUser.Role().IsOwner() {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanSendMessage(conversationUser *entity.ConversationUser) error {
	return nil
}

func (p *ConversationUserPolicy) CanEditMessage(conversationUser *entity.ConversationUser, message *entity.Message) error {
	if !conversationUser.ConversationID().Equal(message.ConversationID()) {
		return errors.New("")
	}

	if message.HasUser() && conversationUser.UserID().Equal(*message.UserID()) {
		return nil
	}

	return errors.New("")
}

func (p *ConversationUserPolicy) CanDeleteMessage(conversationUser *entity.ConversationUser, message *entity.Message) error {
	if !conversationUser.ConversationID().Equal(message.ConversationID()) {
		return errors.New("")
	}

	if message.HasUser() && conversationUser.UserID().Equal(*message.UserID()) {
		return nil
	}
	if conversationUser.Role().IsOwner() || conversationUser.Role().IsAdmin() {
		return nil
	}

	return errors.New("")
}
