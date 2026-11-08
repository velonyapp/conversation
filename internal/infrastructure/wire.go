package infrastructure

import (
	"github.com/velonyapp/conversation/internal/infrastructure/data/mysql"
	"github.com/velonyapp/conversation/internal/infrastructure/messaging/protobuf"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewConversationRepo,
	mysql.NewUserRepo,
	mysql.NewConversationUserRepo,
	mysql.NewMessageRepo,
	mysql.NewConversationQuery,
	mysql.NewConversationUserQuery,
	mysql.NewMessageQuery,
	mysql.NewEventPublisher,
	protobuf.NewEventEncoder,
)
