package repo

import (
	"context"

	"github.com/velonyapp/conversation/internal/domain/entity"
	"github.com/velonyapp/conversation/internal/domain/vo"
)

type User interface {
	FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error)

	Save(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userID vo.UserID) error
}
