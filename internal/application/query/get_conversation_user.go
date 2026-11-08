package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type GetConversationUser struct {
	ActorID string

	ConversationID     string
	ConversationUserID string
}

type GetConversationUserResult struct {
	ConversationUser common.ConversationUserResult
}

func (GetConversationUser) resultType() GetConversationUserResult {
	return GetConversationUserResult{}
}

type GetConversationUserHandler Handler[GetConversationUser, GetConversationUserResult]

type getConversationUserHandler struct {
	conversationUserQuery ConversationUser
}

func NewGetConversationUserHandler(
	conversationUserQuery ConversationUser,
) GetConversationUserHandler {
	return &getConversationUserHandler{
		conversationUserQuery: conversationUserQuery,
	}
}

func (h *getConversationUserHandler) Handle(
	ctx context.Context,
	qry GetConversationUser,
) (GetConversationUserResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)
	conversationUserID, _ := vo.NewConversationUserID(qry.ConversationUserID)

	conversationUser, err := h.conversationUserQuery.Find(ctx,
		conversationUserID,
		ConversationUserScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
	)
	if err != nil {
		return GetConversationUserResult{}, err
	}
	if conversationUser == nil {
		return GetConversationUserResult{}, common.ErrConversationUserNotFound
	}

	return GetConversationUserResult{
		ConversationUser: *conversationUser,
	}, nil
}
