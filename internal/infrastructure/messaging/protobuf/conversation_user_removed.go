package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationUserRemovedPayload(event integrationevent.ConversationUserRemoved) *v1.ConversationUserRemovedPayload {
	return &v1.ConversationUserRemovedPayload{
		ConversationUserId: event.AggregateID(),
		ConversationId:     event.ConversationID(),
		UserId:             event.UserID(),
		RemoveTime:         timestamppb.New(event.OccurTime()),
	}
}
