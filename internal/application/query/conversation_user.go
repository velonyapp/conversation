package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ConversationUserScope struct {
	ConversationID vo.ConversationID
	UserID         vo.UserID
}

type ConversationUser interface {
	Find(ctx context.Context, conversationUserID vo.ConversationUserID, scope ConversationUserScope) (*common.ConversationUserResult, error)
	BatchFind(ctx context.Context, conversationUserIDs []vo.ConversationUserID, scope ConversationUserScope) (map[vo.ConversationUserID]common.ConversationUserResult, error)
	List(ctx context.Context, scope ConversationUserScope, pageSize int32, pageToken string) ([]common.ConversationUserResult, string, error)
}
