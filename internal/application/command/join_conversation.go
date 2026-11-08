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

type JoinConversation struct {
	ActorID string

	ConversationID string
}

type JoinConversationResult struct {
	Conversation common.ConversationResult
}

func (JoinConversation) resultType() JoinConversationResult {
	return JoinConversationResult{}
}

type JoinConversationHandler Handler[JoinConversation, JoinConversationResult]

type joinConversationHandler struct {
	conversationRepo     repo.Conversation
	userRepo             repo.User
	conversationUserRepo repo.ConversationUser
	messageRepo          repo.Message
	unitOfWork           port.UnitOfWork
}

func NewJoinConversationHandler(
	conversationRepo repo.Conversation,
	userRepo repo.User,
	conversationUserRepo repo.ConversationUser,
	messageRepo repo.Message,
	unitOfWork port.UnitOfWork,
) JoinConversationHandler {
	return &joinConversationHandler{
		conversationUserRepo: conversationUserRepo,
		userRepo:             userRepo,
		conversationRepo:     conversationRepo,
		messageRepo:          messageRepo,
		unitOfWork:           unitOfWork,
	}
}

func (h *joinConversationHandler) Handle(
	ctx context.Context,
	cmd JoinConversation,
) (JoinConversationResult, error) {
	now := time.Now().UTC()

	actorID := vo.NewUserID(cmd.ActorID)

	conversationID, _ := vo.NewConversationID(cmd.ConversationID)

	actor, err := h.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return JoinConversationResult{}, err
	}
	if actor == nil {
		return JoinConversationResult{}, common.ErrUserNotFound
	}

	conversation, err := h.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return JoinConversationResult{}, err
	}
	if conversation == nil {
		return JoinConversationResult{}, common.ErrConversationNotFound
	}

	conversationUser := entity.NewConversationUser(
		conversation.ID(),
		actor.ID(),
		vo.ConversationUserRoleMember,
		now,
	)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.conversationUserRepo.Save(ctx, conversationUser)
	}); err != nil {
		return JoinConversationResult{}, err
	}

	conversationResult := common.ConversationResult{
		ID:          conversation.ID().String(),
		Title:       conversation.Title().String(),
		Description: conversation.Description().String(),
		CreateTime:  conversation.CreateTime(),
		UpdateTime:  conversation.UpdateTime(),
	}
	lastMessage, err := h.messageRepo.FindLastByConversationID(ctx, conversationID)
	if err != nil {
		return JoinConversationResult{}, err
	}
	if lastMessage != nil {
		value := lastMessage.ID().String()
		conversationResult.LastMessageID = &value
	}

	return JoinConversationResult{
		Conversation: conversationResult,
	}, nil
}
