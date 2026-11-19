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
	conversationUserIDs := make([]vo.ConversationUserID, 0, len(qry.ConversationUserIDs))
	hasMe := false
	for _, conversationUserID := range qry.ConversationUserIDs {
		if conversationUserID == conversationUserAliasMe {
			hasMe = true
		} else {
			id, _ := vo.NewConversationUserID(conversationUserID)
			conversationUserIDs = append(conversationUserIDs, id)
		}
	}

	scope := ConversationUserScope{
		ConversationID: conversationID,
		UserID:         actorID,
	}

	conversationUsers, err := h.conversationUserQuery.BatchFind(ctx, conversationUserIDs, scope)
	if err != nil {
		return BatchGetConversationUsersResult{}, err
	}

	var me *common.ConversationUserResult
	if hasMe {
		me, err = h.conversationUserQuery.FindMe(ctx, scope)
		if err != nil {
			return BatchGetConversationUsersResult{}, err
		}
		if me == nil {
			return BatchGetConversationUsersResult{}, common.ErrConversationUserNotFound
		}
	}

	conversationUserResults := make([]common.ConversationUserResult, len(qry.ConversationUserIDs))
	for i, conversationUserID := range qry.ConversationUserIDs {
		if conversationUserID == conversationUserAliasMe {
			conversationUserResults[i] = *me
		} else {
			id, _ := vo.NewConversationUserID(conversationUserID)
			conversationUser, ok := conversationUsers[id]
			if !ok {
				return BatchGetConversationUsersResult{}, common.ErrConversationUserNotFound
			}

			conversationUserResults[i] = conversationUser
		}
	}

	return BatchGetConversationUsersResult{
		ConversationUsers: conversationUserResults,
	}, nil
}
