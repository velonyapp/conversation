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

type DeleteConversation struct {
	ActorID string

	ConversationID string
}

type DeleteConversationResult struct {
}

func (DeleteConversation) resultType() DeleteConversationResult {
	return DeleteConversationResult{}
}

type DeleteConversationHandler Handler[DeleteConversation, DeleteConversationResult]

type deleteConversationHandler struct {
	conversationRepo       repo.Conversation
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewDeleteConversationHandler(
	conversationRepo repo.Conversation,
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) DeleteConversationHandler {
	return &deleteConversationHandler{
		conversationRepo:       conversationRepo,
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *deleteConversationHandler) Handle(
	ctx context.Context,
	cmd DeleteConversation,
) (DeleteConversationResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return DeleteConversationResult{}, err
	}
	if actor == nil {
		return DeleteConversationResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanDeleteConversation(actor); err != nil {
		return DeleteConversationResult{}, err
	}

	conversation, err := h.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return DeleteConversationResult{}, err
	}
	if conversation == nil {
		return DeleteConversationResult{}, common.ErrConversationNotFound
	}

	if err := conversation.Delete(now); err != nil {
		return DeleteConversationResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationRepo.Save(ctx, conversation)
	}); err != nil {
		return DeleteConversationResult{}, err
	}

	return DeleteConversationResult{}, nil
}
