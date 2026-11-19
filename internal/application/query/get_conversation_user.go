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
	scope := ConversationUserScope{
		ConversationID: conversationID,
		UserID:         actorID,
	}

	var (
		conversationUser *common.ConversationUserResult
		err              error
	)
	if qry.ConversationUserID == conversationUserAliasMe {
		conversationUser, err = h.conversationUserQuery.FindMe(ctx, scope)
	} else {
		conversationUserID, _ := vo.NewConversationUserID(qry.ConversationUserID)
		conversationUser, err = h.conversationUserQuery.Find(ctx, conversationUserID, scope)
	}
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
