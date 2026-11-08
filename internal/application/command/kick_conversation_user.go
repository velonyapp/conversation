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

type KickConversationUser struct {
	ActorID string

	ConversationID string
	TargetID       string
}

type KickConversationUserResult struct {
}

func (KickConversationUser) resultType() KickConversationUserResult {
	return KickConversationUserResult{}
}

type KickConversationUserHandler Handler[KickConversationUser, KickConversationUserResult]

type kickConversationUserHandler struct {
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewKickConversationUserHandler(
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) KickConversationUserHandler {
	return &kickConversationUserHandler{
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *kickConversationUserHandler) Handle(
	ctx context.Context,
	cmd KickConversationUser,
) (KickConversationUserResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return KickConversationUserResult{}, err
	}
	if actor == nil {
		return KickConversationUserResult{}, common.ErrConversationUserNotFound
	}

	targetID, _ := vo.NewConversationUserID(cmd.TargetID)

	target, err := h.conversationUserRepo.FindByID(ctx, targetID)
	if err != nil {
		return KickConversationUserResult{}, err
	}
	if target == nil {
		return KickConversationUserResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanKickOther(actor, target); err != nil {
		return KickConversationUserResult{}, err
	}

	if err := target.Remove(now); err != nil {
		return KickConversationUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, target)
	}); err != nil {
		return KickConversationUserResult{}, err
	}

	return KickConversationUserResult{}, nil
}
