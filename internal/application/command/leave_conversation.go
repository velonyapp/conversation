package command

import (
	"context"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type LeaveConversation struct {
	ActorID string

	ConversationID string
}

type LeaveConversationResult struct {
}

func (LeaveConversation) resultType() LeaveConversationResult {
	return LeaveConversationResult{}
}

type LeaveConversationHandler Handler[LeaveConversation, LeaveConversationResult]

type leaveConversationHandler struct {
	conversationUserRepo repo.ConversationUser
	unitOfWork           port.UnitOfWork
}

func NewLeaveConversationHandler(
	conversationUserRepo repo.ConversationUser,
	unitOfWork port.UnitOfWork,
) LeaveConversationHandler {
	return &leaveConversationHandler{
		conversationUserRepo: conversationUserRepo,
		unitOfWork:           unitOfWork,
	}
}

func (h *leaveConversationHandler) Handle(
	ctx context.Context,
	cmd LeaveConversation,
) (LeaveConversationResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	conversationUser, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx, conversationID, actorID)
	if err != nil {
		return LeaveConversationResult{}, err
	}
	if conversationUser == nil {
		return LeaveConversationResult{}, common.ErrConversationUserNotFound
	}

	if err := conversationUser.Remove(now); err != nil {
		return LeaveConversationResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, conversationUser)
	}); err != nil {
		return LeaveConversationResult{}, err
	}

	return LeaveConversationResult{}, nil
}
