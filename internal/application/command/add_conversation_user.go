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

type AddConversationUser struct {
	ActorID string

	ConversationID string
	UserID         string
}

type AddConversationUserResult struct {
	ConversationUser common.ConversationUserResult
}

func (AddConversationUser) resultType() AddConversationUserResult {
	return AddConversationUserResult{}
}

type AddConversationUserHandler Handler[AddConversationUser, AddConversationUserResult]

type addConversationUserHandler struct {
	userRepo               repo.User
	conversationUserRepo   repo.ConversationUser
	conversationUserPolicy *service.ConversationUserPolicy
	unitOfWork             port.UnitOfWork
}

func NewAddConversationUserHandler(
	userRepo repo.User,
	conversationUserRepo repo.ConversationUser,
	conversationUserPolicy *service.ConversationUserPolicy,
	unitOfWork port.UnitOfWork,
) AddConversationUserHandler {
	return &addConversationUserHandler{
		userRepo:               userRepo,
		conversationUserRepo:   conversationUserRepo,
		conversationUserPolicy: conversationUserPolicy,
		unitOfWork:             unitOfWork,
	}
}

func (h *addConversationUserHandler) Handle(
	ctx context.Context,
	cmd AddConversationUser,
) (AddConversationUserResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.conversationUserRepo.FindByConversationIDAndUserID(ctx,
		conversationID,
		actorID,
	)
	if err != nil {
		return AddConversationUserResult{}, err
	}
	if actor == nil {
		return AddConversationUserResult{}, common.ErrConversationUserNotFound
	}

	if err := h.conversationUserPolicy.CanAddUser(actor); err != nil {
		return AddConversationUserResult{}, err
	}

	userID := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return AddConversationUserResult{}, err
	}
	if user == nil {
		return AddConversationUserResult{}, common.ErrUserNotFound
	}

	conversationUser := entity.NewConversationUser(
		conversationID,
		user.ID(),
		vo.ConversationUserRoleMember,
		now,
	)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, conversationUser)
	}); err != nil {
		return AddConversationUserResult{}, err
	}

	return AddConversationUserResult{
		ConversationUser: common.ConversationUserResult{
			ID:         conversationUser.ID().String(),
			UserID:     conversationUser.UserID().String(),
			Role:       conversationUser.Role().String(),
			CreateTime: conversationUser.CreateTime(),
			UpdateTime: conversationUser.UpdateTime(),
		},
	}, nil
}
