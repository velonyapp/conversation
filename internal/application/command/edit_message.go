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

type EditMessage struct {
	ActorID string

	ConversationID string
	MessageID      string
	Content        string
}

type EditMessageResult struct {
	Message common.MessageResult
}

func (EditMessage) resultType() EditMessageResult {
	return EditMessageResult{}
}

type EditMessageHandler Handler[EditMessage, EditMessageResult]

type editMessageHandler struct {
	conversationUserRepo   repo.ConversationUser
	messageRepo            repo.Message
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewEditMessageHandler(
	conversationUserRepo repo.ConversationUser,
	messageRepo repo.Message,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) EditMessageHandler {
	return &editMessageHandler{
		conversationUserRepo:   conversationUserRepo,
		messageRepo:            messageRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *editMessageHandler) Handle(
	ctx context.Context,
	cmd EditMessage,
) (EditMessageResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return EditMessageResult{}, err
	}
	if actor == nil {
		return EditMessageResult{}, common.ErrConversationUserNotFound
	}

	messageID, _ := vo.NewMessageID(cmd.MessageID)

	message, err := h.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return EditMessageResult{}, err
	}
	if message == nil {
		return EditMessageResult{}, common.ErrMessageNotFound
	}

	if err := h.conversationUserPolicy.CanEditMessage(actor, message); err != nil {
		return EditMessageResult{}, err
	}

	content, err := vo.NewMessageContent(cmd.Content)
	if err != nil {
		return EditMessageResult{}, err
	}

	if err := message.Edit(content, now); err != nil {
		return EditMessageResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.messageRepo.Save(ctx, message)
	}); err != nil {
		return EditMessageResult{}, err
	}

	messageResult := common.MessageResult{
		ID:         message.ID().String(),
		Content:    message.Content().String(),
		CreateTime: message.CreateTime(),
		UpdateTime: message.UpdateTime(),
	}
	if message.HasUser() {
		value := message.UserID().String()
		messageResult.UserID = &value
	}

	return EditMessageResult{
		Message: messageResult,
	}, nil
}
