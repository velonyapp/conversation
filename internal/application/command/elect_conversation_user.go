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

type ElectConversationUser struct {
	ActorID string

	ConversationID string
	TargetID       string
}

type ElectConversationUserResult struct {
	ConversationUser common.ConversationUserResult
}

func (ElectConversationUser) resultType() ElectConversationUserResult {
	return ElectConversationUserResult{}
}

type ElectConversationUserHandler Handler[ElectConversationUser, ElectConversationUserResult]

type electConversationUserHandler struct {
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewElectConversationUserHandler(
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) ElectConversationUserHandler {
	return &electConversationUserHandler{
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *electConversationUserHandler) Handle(
	ctx context.Context,
	cmd ElectConversationUser,
) (ElectConversationUserResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return ElectConversationUserResult{}, err
	}
	if actor == nil {
		return ElectConversationUserResult{}, common.ErrConversationUserNotFound
	}

	targetID, _ := vo.NewConversationUserID(cmd.TargetID)

	target, err := h.conversationUserRepo.FindByID(ctx, targetID)
	if err != nil {
		return ElectConversationUserResult{}, err
	}
	if target == nil {
		return ElectConversationUserResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanElectOther(actor, target); err != nil {
		return ElectConversationUserResult{}, err
	}

	if err := actor.ChangeRole(vo.ConversationUserRoleAdmin, now); err != nil {
		return ElectConversationUserResult{}, err
	}
	if err := target.ChangeRole(vo.ConversationUserRoleOwner, now); err != nil {
		return ElectConversationUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.conversationUserRepo.Save(ctx, actor); err != nil {
			return err
		}

		return h.conversationUserRepo.Save(ctx, target)
	}); err != nil {
		return ElectConversationUserResult{}, err
	}

	return ElectConversationUserResult{
		ConversationUser: common.ConversationUserResult{
			ID:         target.ID().String(),
			UserID:     target.UserID().String(),
			Role:       target.Role().String(),
			CreateTime: target.CreateTime(),
			UpdateTime: target.UpdateTime(),
		},
	}, nil
}
