package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ConversationScope struct {
	UserID vo.UserID
}

type Conversation interface {
	Find(ctx context.Context, conversationID vo.ConversationID, scope ConversationScope) (*common.ConversationResult, error)
	BatchFind(ctx context.Context, conversationIDs []vo.ConversationID, scope ConversationScope) (map[vo.ConversationID]common.ConversationResult, error)
	List(ctx context.Context, scope ConversationScope, pageSize int32, pageToken string) ([]common.ConversationResult, string, error)
}
