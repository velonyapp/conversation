package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationUpdatedPayload(event integrationevent.ConversationUpdated) *v1.ConversationUpdatedPayload {
	return &v1.ConversationUpdatedPayload{
		ConversationId: event.AggregateID(),
		Title:          event.Title(),
		Description:    event.Description(),
		UpdateTime:     timestamppb.New(event.OccurTime()),
		Version:        event.Version(),
	}
}
