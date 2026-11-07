package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrConversationUserVersionConflict = errors.New("conversation user version conflict")
)

type ConversationUser interface {
	FindByID(ctx context.Context, conversationUserID vo.ConversationUserID) (*entity.ConversationUser, error)
	FindByConversationIDAndUserID(ctx context.Context, conversationID vo.ConversationID, userID vo.UserID) (*entity.ConversationUser, error)

	Save(ctx context.Context, conversation *entity.ConversationUser) error
}
