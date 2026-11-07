package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ListConversations struct {
	ActorID string

	PageSize  int32
	PageToken string
}

type ListConversationsResult struct {
	Conversations []common.ConversationResult
	NextPageToken string
}

func (ListConversations) resultType() ListConversationsResult {
	return ListConversationsResult{}
}

type ListConversationsHandler Handler[ListConversations, ListConversationsResult]

type listConversationsHandler struct {
	conversationQuery Conversation
}

func NewListConversationsHandler(
	conversationQuery Conversation,
) ListConversationsHandler {
	return &listConversationsHandler{
		conversationQuery: conversationQuery,
	}
}

func (h *listConversationsHandler) Handle(
	ctx context.Context,
	qry ListConversations,
) (ListConversationsResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversations, nextPageToken, err := h.conversationQuery.List(ctx,
		ConversationScope{
			UserID: actorID,
		},
		qry.PageSize,
		qry.PageToken,
	)
	if err != nil {
		return ListConversationsResult{}, err
	}

	return ListConversationsResult{
		Conversations: conversations,
		NextPageToken: nextPageToken,
	}, nil
}
