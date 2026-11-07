package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type BatchGetConversations struct {
	ActorID string

	ConversationIDs []string
}

type BatchGetConversationsResult struct {
	Conversations []common.ConversationResult
}

func (BatchGetConversations) resultType() BatchGetConversationsResult {
	return BatchGetConversationsResult{}
}

type BatchGetConversationsHandler Handler[BatchGetConversations, BatchGetConversationsResult]

type batchGetConversationsHandler struct {
	conversationQuery Conversation
}

func NewBatchGetConversationsHandler(
	conversationQuery Conversation,
) BatchGetConversationsHandler {
	return &batchGetConversationsHandler{
		conversationQuery: conversationQuery,
	}
}

func (h *batchGetConversationsHandler) Handle(
	ctx context.Context,
	qry BatchGetConversations,
) (BatchGetConversationsResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationIDs := make([]vo.ConversationID, len(qry.ConversationIDs))
	for i, conversationID := range qry.ConversationIDs {
		conversationIDs[i], _ = vo.NewConversationID(conversationID)
	}

	conversations, err := h.conversationQuery.BatchFind(ctx,
		conversationIDs,
		ConversationScope{
			UserID: actorID,
		},
	)
	if err != nil {
		return BatchGetConversationsResult{}, err
	}

	conversationResults := make([]common.ConversationResult, len(conversationIDs))
	for i, conversationID := range conversationIDs {
		conversation, ok := conversations[conversationID]
		if !ok {
			return BatchGetConversationsResult{}, common.ErrConversationNotFound
		}

		conversationResults[i] = conversation
	}

	return BatchGetConversationsResult{
		Conversations: conversationResults,
	}, nil
}
