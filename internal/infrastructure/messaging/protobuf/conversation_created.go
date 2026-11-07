package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationCreatedPayload(event integrationevent.ConversationCreated) *v1.ConversationCreatedPayload {
	return &v1.ConversationCreatedPayload{
		ConversationId: event.AggregateID(),
		Title:          event.Title(),
		Description:    event.Description(),
		CreateTime:     timestamppb.New(event.OccurTime()),
		Version:        event.Version(),
	}
}
