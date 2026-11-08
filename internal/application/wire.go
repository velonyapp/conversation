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
	command.NewJoinConversationHandler,
	command.NewLeaveConversationHandler,
	command.NewAddConversationUserHandler,
	command.NewKickConversationUserHandler,
	command.NewElectConversationUserHandler,
	command.NewPromoteConversationUserHandler,
	command.NewDemoteConversationUserHandler,
	command.NewCreateMessageHandler,
	command.NewEditMessageHandler,
	command.NewDeleteMessageHandler,
	command.NewCreateUserHandler,
	command.NewDeleteUserHandler,
	query.NewBus,
	query.NewGetConversationHandler,
	query.NewListConversationsHandler,
	query.NewBatchGetConversationsHandler,
	query.NewGetConversationUserHandler,
	query.NewListConversationUsersHandler,
	query.NewBatchGetConversationUsersHandler,
	query.NewGetMessageHandler,
	query.NewListMessagesHandler,
	query.NewBatchGetMessagesHandler,
	domainevent.NewDispatcher,
	domainevent.NewConversationCreatedHandler,
	domainevent.NewConversationUpdatedHandler,
	domainevent.NewConversationDeletedHandler,
	domainevent.NewConversationUserAddedHandler,
	domainevent.NewConversationUserUpdatedHandler,
	domainevent.NewConversationUserRemovedHandler,
	domainevent.NewMessageCreatedHandler,
	domainevent.NewMessageEditedHandler,
	domainevent.NewMessageDeletedHandler,
)
