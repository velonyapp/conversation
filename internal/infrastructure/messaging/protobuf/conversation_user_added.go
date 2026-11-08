package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationUserAddedPayload(event integrationevent.ConversationUserAdded) *v1.ConversationUserAddedPayload {
	return &v1.ConversationUserAddedPayload{
		ConversationUserId: event.AggregateID(),
		ConversationId:     event.ConversationID(),
		UserId:             event.UserID(),
		Role:               event.Role(),
		AddTime:            timestamppb.New(event.OccurTime()),
		Version:            event.Version(),
	}
}
