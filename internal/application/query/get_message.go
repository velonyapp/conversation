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

	messageID, _ := vo.NewMessageID(qry.MessageID)
	conversationID, _ := vo.NewConversationID(qry.ConversationID)

	message, err := h.messageQuery.Find(ctx,
		messageID,
		MessageScope{
			ConversationID: conversationID,
			UserID:         actorID,
		},
	)
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
