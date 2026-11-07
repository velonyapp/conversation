package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationUserUpdatedPayload(event integrationevent.ConversationUserUpdated) *v1.ConversationUserUpdatedPayload {
	return &v1.ConversationUserUpdatedPayload{
		ConversationUserId: event.AggregateID(),
		ConversationId:     event.ConversationID(),
		UserId:             event.UserID(),
		Role:               event.Role(),
		UpdateTime:         timestamppb.New(event.OccurTime()),
		Version:            event.Version(),
	}
}
