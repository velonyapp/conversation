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

type UpdateConversation struct {
	ActorID string

	ConversationID string
	Title          *string
	Description    *string
}

type UpdateConversationResult struct {
	Conversation common.ConversationResult
}

func (UpdateConversation) resultType() UpdateConversationResult {
	return UpdateConversationResult{}
}

type UpdateConversationHandler Handler[UpdateConversation, UpdateConversationResult]

type updateConversationHandler struct {
	conversationRepo       repo.Conversation
	conversationUserRepo   repo.ConversationUser
	messageRepo            repo.Message
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewUpdateConversationHandler(
	conversationRepo repo.Conversation,
	conversationUserRepo repo.ConversationUser,
	messageRepo repo.Message,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) UpdateConversationHandler {
	return &updateConversationHandler{
		conversationRepo:       conversationRepo,
		conversationUserRepo:   conversationUserRepo,
		messageRepo:            messageRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *updateConversationHandler) Handle(
	ctx context.Context,
	cmd UpdateConversation,
) (UpdateConversationResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	conversationUser, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return UpdateConversationResult{}, err
	}
	if conversationUser == nil {
		return UpdateConversationResult{}, common.ErrConversationUserNotFound
	}

	conversation, err := h.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return UpdateConversationResult{}, err
	}
	if conversation == nil {
		return UpdateConversationResult{}, common.ErrConversationNotFound
	}

	if err := h.conversationUserPolicy.CanUpdateConversation(conversationUser, conversation); err != nil {
		return UpdateConversationResult{}, err
	}

	if cmd.Title != nil {
		title, err := vo.NewConversationTitle(*cmd.Title)
		if err != nil {
			return UpdateConversationResult{}, err
		}

		if err := conversation.ChangeTitle(title, now); err != nil {
			return UpdateConversationResult{}, err
		}
	}
	if cmd.Description != nil {
		title, err := vo.NewConversationDescription(*cmd.Description)
		if err != nil {
			return UpdateConversationResult{}, err
		}

		if err := conversation.ChangeDescription(title, now); err != nil {
			return UpdateConversationResult{}, err
		}
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationRepo.Save(ctx, conversation)
	}); err != nil {
		return UpdateConversationResult{}, err
	}

	conversationResult := common.ConversationResult{
		ID:          conversation.ID().String(),
		Title:       conversation.Title().String(),
		Description: conversation.DeleteTime().String(),
		CreateTime:  conversation.CreateTime(),
		UpdateTime:  conversation.UpdateTime(),
	}
	lastMessage, err := h.messageRepo.FindLastByConversationID(ctx, conversationID)
	if err != nil {
		return UpdateConversationResult{}, err
	}
	if lastMessage != nil {
		value := lastMessage.ID().String()
		conversationResult.LastMessageID = &value
	}

	return UpdateConversationResult{
		Conversation: conversationResult,
	}, nil
}
