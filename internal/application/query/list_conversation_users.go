package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ListConversationUsers struct {
	ActorID string

	ConversationID string
	PageSize       int32
	PageToken      string
}

type ListConversationUsersResult struct {
	ConversationUsers []common.ConversationUserResult
	NextPageToken     string
}

func (ListConversationUsers) resultType() ListConversationUsersResult {
	return ListConversationUsersResult{}
}

type ListConversationUsersHandler Handler[ListConversationUsers, ListConversationUsersResult]

type listConversationUsersHandler struct {
	conversationUserQuery ConversationUser
}

func NewListConversationUsersHandler(
	conversationUserQuery ConversationUser,
) ListConversationUsersHandler {
	return &listConversationUsersHandler{
		conversationUserQuery: conversationUserQuery,
	}
}

func (h *listConversationUsersHandler) Handle(
	ctx context.Context,
	qry ListConversationUsers,
) (ListConversationUsersResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)

	conversationUsers, nextPageToken, err := h.conversationUserQuery.List(ctx,
		ConversationUserScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
		qry.PageSize,
		qry.PageToken,
	)
	if err != nil {
		return ListConversationUsersResult{}, err
	}

	return ListConversationUsersResult{
		ConversationUsers: conversationUsers,
		NextPageToken:     nextPageToken,
	}, nil
}
