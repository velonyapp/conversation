package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func messageDeletedPayload(event integrationevent.MessageDeleted) *v1.MessageDeletedPayload {
	return &v1.MessageDeletedPayload{
		MessageId:      event.AggregateID(),
		ConversationId: event.ConversationID(),
		UserId:         event.UserID(),
		DeleteTime:     timestamppb.New(event.OccurTime()),
	}
}
