package query

import (
	"context"

	"github.com/velonyapp/conversation/internal/application/common"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type ConversationUserScope struct {
	ConversationID vo.ConversationID
}

type ConversationUser interface {
	Get(ctx context.Context, userID vo.UserID, scope ConversationUserScope) (*common.ConversationUserResult, error)
	BatchGet(ctx context.Context, userIDs []vo.UserID, scope ConversationUserScope) (map[vo.ConversationUserID]common.ConversationUserResult, error)
	List(ctx context.Context, scope ConversationUserScope, pageSize int32, pageToken string) ([]common.ConversationUserResult, string, error)
}
