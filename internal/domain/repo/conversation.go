package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrConversationVersionConflict = errors.New("conversation version conflict")
)

type Conversation interface {
	FindByID(ctx context.Context, conversationID vo.ConversationID) (*entity.Conversation, error)

	Save(ctx context.Context, conversation *entity.Conversation) error
}
