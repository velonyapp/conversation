package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

var (
	ErrMessageVersionConflict = errors.New("message version conflict")
)

type Message interface {
	FindByID(ctx context.Context, messageID vo.MessageID) (*entity.Message, error)
	FindLastByConversationID(ctx context.Context, conversationID vo.ConversationID) (*entity.Message, error)

	Save(ctx context.Context, message *entity.Message) error
}
