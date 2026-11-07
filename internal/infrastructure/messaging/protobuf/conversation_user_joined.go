package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationUserJoinedPayload(event integrationevent.ConversationUserJoined) *v1.ConversationUserJoinedPayload {
	return &v1.ConversationUserJoinedPayload{
		ConversationUserId: event.AggregateID(),
		ConversationId:     event.ConversationID(),
		UserId:             event.UserID(),
		Role:               event.Role(),
		JoinTime:           timestamppb.New(event.OccurTime()),
		Version:            event.Version(),
	}
}
