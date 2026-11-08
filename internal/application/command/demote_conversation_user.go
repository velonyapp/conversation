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

type DemoteConversationUser struct {
	ActorID string

	ConversationID string
	TargetID       string
}

type DemoteConversationUserResult struct {
	ConversationUser common.ConversationUserResult
}

func (DemoteConversationUser) resultType() DemoteConversationUserResult {
	return DemoteConversationUserResult{}
}

type DemoteConversationUserHandler Handler[DemoteConversationUser, DemoteConversationUserResult]

type demoteConversationUserHandler struct {
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewDemoteConversationUserHandler(
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) DemoteConversationUserHandler {
	return &demoteConversationUserHandler{
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *demoteConversationUserHandler) Handle(
	ctx context.Context,
	cmd DemoteConversationUser,
) (DemoteConversationUserResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return DemoteConversationUserResult{}, err
	}
	if actor == nil {
		return DemoteConversationUserResult{}, common.ErrConversationUserNotFound
	}

	targetID, _ := vo.NewConversationUserID(cmd.TargetID)

	target, err := h.conversationUserRepo.FindByID(ctx, targetID)
	if err != nil {
		return DemoteConversationUserResult{}, err
	}
	if target == nil {
		return DemoteConversationUserResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanDemoteOther(actor, target); err != nil {
		return DemoteConversationUserResult{}, err
	}

	if err := target.ChangeRole(vo.ConversationUserRoleMember, now); err != nil {
		return DemoteConversationUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, target)
	}); err != nil {
		return DemoteConversationUserResult{}, err
	}

	return DemoteConversationUserResult{
		ConversationUser: common.ConversationUserResult{
			ID:         target.ID().String(),
			UserID:     target.UserID().String(),
			Role:       target.Role().String(),
			CreateTime: target.CreateTime(),
			UpdateTime: target.UpdateTime(),
		},
	}, nil
}
