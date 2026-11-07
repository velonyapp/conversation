package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type MessageScope struct {
	ConversationID vo.ConversationID
	UserID         vo.UserID
}

type Message interface {
	Find(ctx context.Context, messageID vo.MessageID, scope MessageScope) (*common.MessageResult, error)
	BatchFind(ctx context.Context, messageIDs []vo.MessageID, scope MessageScope) (map[vo.MessageID]common.MessageResult, error)
	List(ctx context.Context, scope MessageScope, pageSize int32, pageToken string) ([]common.MessageResult, string, error)
}
