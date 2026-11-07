package command

import (
	"context"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/service"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type CreateMessage struct {
	ActorID string

	ConversationID string
	Content        string
}

type CreateMessageResult struct {
	Message common.MessageResult
}

func (CreateMessage) resultType() CreateMessageResult {
	return CreateMessageResult{}
}

type CreateMessageHandler Handler[CreateMessage, CreateMessageResult]

type createMessageHandler struct {
	conversationUserRepo   repo.ConversationUser
	messageRepo            repo.Message
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewCreateMessageHandler(
	conversationUserRepo repo.ConversationUser,
	messageRepo repo.Message,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) CreateMessageHandler {
	return &createMessageHandler{
		conversationUserRepo:   conversationUserRepo,
		messageRepo:            messageRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *createMessageHandler) Handle(
	ctx context.Context,
	cmd CreateMessage,
) (CreateMessageResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	conversationUser, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return CreateMessageResult{}, err
	}
	if conversationUser == nil {
		return CreateMessageResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanSendMessage(conversationUser); err != nil {
		return CreateMessageResult{}, err
	}

	content, err := vo.NewMessageContent(cmd.Content)
	if err != nil {
		return CreateMessageResult{}, err
	}

	message := entity.NewMessage(
		conversationID,
		&actorID,
		content,
		now,
	)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.messageRepo.Save(ctx, message)
	}); err != nil {
		return CreateMessageResult{}, err
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

	return CreateMessageResult{
		Message: messageResult,
	}, nil
}
