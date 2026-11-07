package domain

import (
	"github.com/velonyapp/conversation/internal/domain/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewConversationUserPolicy,
)
