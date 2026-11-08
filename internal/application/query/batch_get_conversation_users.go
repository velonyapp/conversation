package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type BatchGetConversationUsers struct {
	ActorID string

	ConversationID      string
	ConversationUserIDs []string
}

type BatchGetConversationUsersResult struct {
	ConversationUsers []common.ConversationUserResult
}

func (BatchGetConversationUsers) resultType() BatchGetConversationUsersResult {
	return BatchGetConversationUsersResult{}
}

type BatchGetConversationUsersHandler Handler[BatchGetConversationUsers, BatchGetConversationUsersResult]

type batchGetConversationUsersHandler struct {
	conversationUserQuery ConversationUser
}

func NewBatchGetConversationUsersHandler(
	conversationUserQuery ConversationUser,
) BatchGetConversationUsersHandler {
	return &batchGetConversationUsersHandler{
		conversationUserQuery: conversationUserQuery,
	}
}

func (h *batchGetConversationUsersHandler) Handle(
	ctx context.Context,
	qry BatchGetConversationUsers,
) (BatchGetConversationUsersResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)
	conversationUserIDs := make([]vo.ConversationUserID, len(qry.ConversationUserIDs))
	for i, conversationUserID := range qry.ConversationUserIDs {
		conversationUserIDs[i], _ = vo.NewConversationUserID(conversationUserID)
	}

	conversationUsers, err := h.conversationUserQuery.BatchFind(ctx,
		conversationUserIDs,
		ConversationUserScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
	)
	if err != nil {
		return BatchGetConversationUsersResult{}, err
	}

	conversationUserResults := make([]common.ConversationUserResult, len(conversationUserIDs))
	for i, conversationUserID := range conversationUserIDs {
		conversationUser, ok := conversationUsers[conversationUserID]
		if !ok {
			return BatchGetConversationUsersResult{}, common.ErrConversationUserNotFound
		}

		conversationUserResults[i] = conversationUser
	}

	return BatchGetConversationUsersResult{
		ConversationUsers: conversationUserResults,
	}, nil
}
