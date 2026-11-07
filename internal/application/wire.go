package application

import (
	"github.com/velonyapp/conversation/internal/application/command"
	"github.com/velonyapp/conversation/internal/application/domainevent"
	"github.com/velonyapp/conversation/internal/application/query"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	command.NewBus,
	command.NewCreateConversationHandler,
	command.NewUpdateConversationHandler,
	command.NewDeleteConversationHandler,
	command.NewCreateMessageHandler,
	command.NewEditMessageHandler,
	command.NewDeleteMessageHandler,
	command.NewCreateUserHandler,
	command.NewDeleteUserHandler,
	query.NewBus,
	query.NewGetConversationHandler,
	query.NewBatchGetConversationsHandler,
	query.NewListConversationsHandler,
	query.NewGetMessageHandler,
	query.NewBatchGetMessagesHandler,
	query.NewListMessagesHandler,
	domainevent.NewDispatcher,
	domainevent.NewConversationCreatedHandler,
	domainevent.NewConversationUpdatedHandler,
	domainevent.NewConversationDeletedHandler,
	domainevent.NewConversationUserJoinedHandler,
	domainevent.NewConversationUserUpdatedHandler,
	domainevent.NewConversationUserRemovedHandler,
	domainevent.NewMessageCreatedHandler,
	domainevent.NewMessageEditedHandler,
	domainevent.NewMessageDeletedHandler,
)
