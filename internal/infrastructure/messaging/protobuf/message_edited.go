package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func messageEditedPayload(event integrationevent.MessageEdited) *v1.MessageEditedPayload {
	return &v1.MessageEditedPayload{
		MessageId:      event.AggregateID(),
		ConversationId: event.ConversationID(),
		UserId:         event.UserID(),
		Content:        event.Content(),
		EditTime:       timestamppb.New(event.OccurTime()),
		Version:        event.Version(),
	}
}
