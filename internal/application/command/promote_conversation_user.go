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

type PromoteConversationUser struct {
	ActorID string

	ConversationID string
	TargetID       string
}

type PromoteConversationUserResult struct {
}

func (PromoteConversationUser) resultType() PromoteConversationUserResult {
	return PromoteConversationUserResult{}
}

type PromoteConversationUserHandler Handler[PromoteConversationUser, PromoteConversationUserResult]

type promoteConversationUserHandler struct {
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewPromoteConversationUserHandler(
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) PromoteConversationUserHandler {
	return &promoteConversationUserHandler{
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *promoteConversationUserHandler) Handle(
	ctx context.Context,
	cmd PromoteConversationUser,
) (PromoteConversationUserResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return PromoteConversationUserResult{}, err
	}
	if actor == nil {
		return PromoteConversationUserResult{}, common.ErrConversationUserNotFound
	}

	targetID, _ := vo.NewConversationUserID(cmd.TargetID)

	target, err := h.conversationUserRepo.FindByID(ctx, targetID)
	if err != nil {
		return PromoteConversationUserResult{}, err
	}
	if target == nil {
		return PromoteConversationUserResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanPromoteOther(actor, target); err != nil {
		return PromoteConversationUserResult{}, err
	}

	if err := target.ChangeRole(vo.ConversationUserRoleAdmin, now); err != nil {
		return PromoteConversationUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, target)
	}); err != nil {
		return PromoteConversationUserResult{}, err
	}

	return PromoteConversationUserResult{}, nil
}
