package command

import (
	"context"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/service"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type DeleteMessage struct {
	ActorID string

	ConversationID string
	MessageID      string
}

type DeleteMessageResult struct {
}

func (DeleteMessage) resultType() DeleteMessageResult {
	return DeleteMessageResult{}
}

type DeleteMessageHandler Handler[DeleteMessage, DeleteMessageResult]

type deleteMessageHandler struct {
	conversationUserRepo   repo.ConversationUser
	messageRepo            repo.Message
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewDeleteMessageHandler(
	conversationUserRepo repo.ConversationUser,
	messageRepo repo.Message,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) DeleteMessageHandler {
	return &deleteMessageHandler{
		conversationUserRepo:   conversationUserRepo,
		messageRepo:            messageRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *deleteMessageHandler) Handle(
	ctx context.Context,
	cmd DeleteMessage,
) (DeleteMessageResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	conversationUser, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return DeleteMessageResult{}, err
	}
	if conversationUser == nil {
		return DeleteMessageResult{}, common.ErrConversationUserNotFound
	}

	messageID, _ := vo.NewMessageID(cmd.MessageID)

	message, err := h.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return DeleteMessageResult{}, err
	}

	if err := h.conversationUserPolicy.CanDeleteMessage(conversationUser, message); err != nil {
		return DeleteMessageResult{}, err
	}

	if err := message.Delete(now); err != nil {
		return DeleteMessageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.messageRepo.Save(ctx, message)
	}); err != nil {
		return DeleteMessageResult{}, err
	}

	return DeleteMessageResult{}, nil
}
