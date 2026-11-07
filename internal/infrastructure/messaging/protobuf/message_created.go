package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func messageCreatedPayload(event integrationevent.MessageCreated) *v1.MessageCreatedPayload {
	return &v1.MessageCreatedPayload{
		MessageId:      event.AggregateID(),
		ConversationId: event.ConversationID(),
		UserId:         event.UserID(),
		Content:        event.Content(),
		CreateTime:     timestamppb.New(event.OccurTime()),
		Version:        event.Version(),
	}
}
