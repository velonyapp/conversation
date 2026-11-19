package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type BatchGetMessages struct {
	ActorID string

	ConversationID string
	Messages       []BatchGetMessage
}

type BatchGetMessage struct {
	ConversationID string
	MessageID      string
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

	messageIDs := make([]vo.MessageID, 0, len(qry.Messages))
	latestConversationIDs := make([]vo.ConversationID, 0, len(qry.Messages))
	for _, message := range qry.Messages {
		if message.MessageID == messageAliasLatest {
			conversationID, _ := vo.NewConversationID(message.ConversationID)
			latestConversationIDs = append(latestConversationIDs, conversationID)
		} else {
			messageID, _ := vo.NewMessageID(message.MessageID)
			messageIDs = append(messageIDs, messageID)
		}
	}

	scope := MessageBatchScope{
		ConversationID: conversationID,
		UserID:         actorID,
	}

	messages, err := h.messageQuery.BatchFind(ctx, messageIDs, scope)
	if err != nil {
		return BatchGetMessagesResult{}, err
	}
	latestMessages, err := h.messageQuery.BatchFindLatest(ctx, latestConversationIDs, scope)
	if err != nil {
		return BatchGetMessagesResult{}, err
	}

	messageResults := make([]common.MessageResult, len(qry.Messages))
	for i, requestedMessage := range qry.Messages {
		var (
			message common.MessageResult
			ok      bool
		)
		if requestedMessage.MessageID == messageAliasLatest {
			conversationID, _ := vo.NewConversationID(requestedMessage.ConversationID)
			message, ok = latestMessages[conversationID]
		} else {
			messageID, _ := vo.NewMessageID(requestedMessage.MessageID)
			message, ok = messages[messageID]
		}
		if !ok {
			return BatchGetMessagesResult{}, common.ErrMessageNotFound
		}

		messageResults[i] = message
	}

	return BatchGetMessagesResult{
		Messages: messageResults,
	}, nil
}
