package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ListMessages struct {
	ActorID string

	ConversationID string
	PageSize       int32
	PageToken      string
}

type ListMessagesResult struct {
	Messages      []common.MessageResult
	NextPageToken string
}

func (ListMessages) resultType() ListMessagesResult {
	return ListMessagesResult{}
}

type ListMessagesHandler Handler[ListMessages, ListMessagesResult]

type listMessagesHandler struct {
	messageQuery Message
}

func NewListMessagesHandler(
	messageQuery Message,
) ListMessagesHandler {
	return &listMessagesHandler{
		messageQuery: messageQuery,
	}
}

func (h *listMessagesHandler) Handle(
	ctx context.Context,
	qry ListMessages,
) (ListMessagesResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)

	messages, nextPageToken, err := h.messageQuery.List(ctx,
		MessageScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
		qry.PageSize,
		qry.PageToken,
	)
	if err != nil {
		return ListMessagesResult{}, err
	}

	return ListMessagesResult{
		Messages:      messages,
		NextPageToken: nextPageToken,
	}, nil
}
