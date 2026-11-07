package protobuf

import (
	v1 "github.com/velonyapp/conversation/gen/event/v1"
	"github.com/velonyapp/conversation/internal/application/integrationevent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func conversationDeletedPayload(event integrationevent.ConversationDeleted) *v1.ConversationDeletedPayload {
	return &v1.ConversationDeletedPayload{
		ConversationId: event.AggregateID(),
		DeleteTime:     timestamppb.New(event.OccurTime()),
	}
}
