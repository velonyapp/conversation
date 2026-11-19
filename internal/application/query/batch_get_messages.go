package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type BatchGetMessages struct {
	ActorID string

	ConversationID string
	MessageIDs     []string
}

type BatchGetMessagesResult struct {
	Messages []common.MessageResult
}

func (BatchGetMessages) resultType() BatchGetMessagesResult {
	return BatchGetMessagesResult{}
}

type BatchGetMessagesHandler Handler[BatchGetMessages, BatchGetMessagesResult]

type batchGetMessagesHandler struct {
	messageQuery Message
}

func NewBatchGetMessagesHandler(
	messageQuery Message,
) BatchGetMessagesHandler {
	return &batchGetMessagesHandler{
		messageQuery: messageQuery,
	}
}

func (h *batchGetMessagesHandler) Handle(
	ctx context.Context,
	qry BatchGetMessages,
) (BatchGetMessagesResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	var conversationID *vo.ConversationID
	if qry.ConversationID != "-" {
		value, _ := vo.NewConversationID(qry.ConversationID)
		conversationID = &value
	}

	messageIDs := make([]vo.MessageID, len(qry.MessageIDs))
	for i, messageID := range qry.MessageIDs {
		messageIDs[i], _ = vo.NewMessageID(messageID)
	}

	messages, err := h.messageQuery.BatchFind(ctx,
		messageIDs,
		MessageBatchScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
	)
	if err != nil {
		return BatchGetMessagesResult{}, err
	}

	messageResults := make([]common.MessageResult, len(messageIDs))
	for i, messageID := range messageIDs {
		message, ok := messages[messageID]
		if !ok {
			return BatchGetMessagesResult{}, common.ErrMessageNotFound
		}

		messageResults[i] = message
	}

	return BatchGetMessagesResult{
		Messages: messageResults,
	}, nil
}
