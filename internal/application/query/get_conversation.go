package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type GetConversation struct {
	ActorID string

	ConversationID string
}

type GetConversationResult struct {
	Conversation common.ConversationResult
}

func (GetConversation) resultType() GetConversationResult {
	return GetConversationResult{}
}

type GetConversationHandler Handler[GetConversation, GetConversationResult]

type getConversationHandler struct {
	conversationQuery Conversation
}

func NewGetConversationHandler(
	conversationQuery Conversation,
) GetConversationHandler {
	return &getConversationHandler{
		conversationQuery: conversationQuery,
	}
}

func (h *getConversationHandler) Handle(
	ctx context.Context,
	qry GetConversation,
) (GetConversationResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)

	conversation, err := h.conversationQuery.Find(ctx,
		conversationID,
		ConversationScope{
			UserID: actorID,
		},
	)
	if err != nil {
		return GetConversationResult{}, err
	}
	if conversation == nil {
		return GetConversationResult{}, common.ErrConversationNotFound
	}

	return GetConversationResult{
		Conversation: *conversation,
	}, nil
}
