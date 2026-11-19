package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type GetMessage struct {
	ActorID string

	ConversationID string
	MessageID      string
}

type GetMessageResult struct {
	Message common.MessageResult
}

func (GetMessage) resultType() GetMessageResult {
	return GetMessageResult{}
}

type GetMessageHandler Handler[GetMessage, GetMessageResult]

type getMessageHandler struct {
	messageQuery Message
}

func NewGetMessageHandler(
	messageQuery Message,
) GetMessageHandler {
	return &getMessageHandler{
		messageQuery: messageQuery,
	}
}

func (h *getMessageHandler) Handle(
	ctx context.Context,
	qry GetMessage,
) (GetMessageResult, error) {
	actorID := vo.NewUserID(qry.ActorID)

	conversationID, _ := vo.NewConversationID(qry.ConversationID)
	scope := MessageScope{
		ConversationID: conversationID,
		UserID:         actorID,
	}

	var (
		message *common.MessageResult
		err     error
	)
	if qry.MessageID == messageAliasLatest {
		message, err = h.messageQuery.FindLatest(ctx, scope)
	} else {
		messageID, _ := vo.NewMessageID(qry.MessageID)
		message, err = h.messageQuery.Find(ctx, messageID, scope)
	}
	if err != nil {
		return GetMessageResult{}, err
	}
	if message == nil {
		return GetMessageResult{}, common.ErrMessageNotFound
	}

	return GetMessageResult{
		Message: *message,
	}, nil
}
