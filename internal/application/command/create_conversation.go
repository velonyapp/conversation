package command

import (
	"context"
	"time"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/application/port"
	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/repo"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type CreateConversation struct {
	ActorID string

	Title       string
	Description string
}

type CreateConversationResult struct {
	Conversation common.ConversationResult
}

func (CreateConversation) resultType() CreateConversationResult {
	return CreateConversationResult{}
}

type CreateConversationHandler Handler[CreateConversation, CreateConversationResult]

type createConversationHandler struct {
	conversationRepo     repo.Conversation
	userRepo             repo.User
	conversationUserRepo repo.ConversationUser
	unitOfWork           port.UnitOfWork
}

func NewCreateConversationHandler(
	conversationRepo repo.Conversation,
	userRepo repo.User,
	conversationUserRepo repo.ConversationUser,
	unitOfWork port.UnitOfWork,
) CreateConversationHandler {
	return &createConversationHandler{
		conversationRepo:     conversationRepo,
		userRepo:             userRepo,
		conversationUserRepo: conversationUserRepo,
		unitOfWork:           unitOfWork,
	}
}

func (h *createConversationHandler) Handle(
	ctx context.Context,
	cmd CreateConversation,
) (CreateConversationResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	actor, err := h.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return CreateConversationResult{}, err
	}
	if actor == nil {
		return CreateConversationResult{}, common.ErrUserNotFound
	}

	title, err := vo.NewConversationTitle(cmd.Title)
	if err != nil {
		return CreateConversationResult{}, err
	}
	description, err := vo.NewConversationDescription(cmd.Description)
	if err != nil {
		return CreateConversationResult{}, err
	}

	conversation := entity.NewConversation(
		title,
		description,
		now,
	)

	conversationUser := entity.NewConversationUser(
		conversation.ID(),
		actor.ID(),
		vo.ConversationUserRoleOwner,
		now,
	)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.conversationRepo.Save(ctx, conversation); err != nil {
			return err
		}

		return h.conversationUserRepo.Save(ctx, conversationUser)
	}); err != nil {
		return CreateConversationResult{}, err
	}

	return CreateConversationResult{
		Conversation: common.ConversationResult{
			ID:          conversation.ID().String(),
			Title:       conversation.Title().String(),
			Description: conversation.Description().String(),
			CreateTime:  conversation.CreateTime(),
			UpdateTime:  conversation.UpdateTime(),
		},
	}, nil
}
