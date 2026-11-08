package transport

import (
	"context"
	"errors"
	"strings"

	"github.com/velonyapp/conversation/internal/conf"
	"github.com/velonyapp/conversation/internal/presentation/service"
	identityv1 "github.com/velonyapp/identity/gen/event/v1"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v3/transport"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

const (
	groupID = "identity"

	identityUserTopic = "identity.user"

	identityUserCreatedEvent = "identity.user.created"
	identityUserDeletedEvent = "identity.user.deleted"
)

var _ transport.Server = (*KafkaConsumer)(nil)

type KafkaConsumer struct {
	c   *conf.Transport
	svc *service.Consumer

	consumer *kafka.Consumer
}

func NewKafkaConsumer(
	c *conf.Transport,
	svc *service.Consumer,
) *KafkaConsumer {
	return &KafkaConsumer{
		c:   c,
		svc: svc,
	}
}

func (k *KafkaConsumer) handleMessage(ctx context.Context, message *kafka.Message) error {
	ctx, cancel := context.WithTimeout(ctx, k.c.Kafka.Timeout.AsDuration())
	defer cancel()

	var eventType string
	for _, header := range message.Headers {
		switch header.Key {
		case "type":
			eventType = string(header.Value)
		}
	}

	switch eventType {
	case identityUserCreatedEvent:
		event := new(identityv1.UserCreatedPayload)
		if err := protojson.Unmarshal(message.Value, event); err != nil {
			return ErrInvalidMessage
		}

		if err := k.svc.CreateUser(ctx, service.CreateUserRequest{
			UserID: event.UserId,
		}); err != nil {
			return err
		}

	case identityUserDeletedEvent:
		event := new(identityv1.UserDeletedPayload)
		if err := protojson.Unmarshal(message.Value, event); err != nil {
			return ErrInvalidMessage
		}

		if err := k.svc.DeleteUser(ctx, service.DeleteUserRequest{
			UserID: event.UserId,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaConsumer) Start(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(k.c.Kafka.Brokers, ","),
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	})
	if err != nil {
		return err
	}

	k.consumer = consumer

	topics := []string{
		identityUserTopic,
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	for !consumer.IsClosed() {
		event := consumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			// TODO: apparantly poll fetches next message instead of retrying the failed one
			if err := k.handleMessage(ctx, e); err != nil && !errors.Is(err, ErrInvalidMessage) {
				break
			}

			consumer.StoreMessage(e)

		case kafka.Error:
			if e.IsFatal() {
				return errors.New(e.Error())
			}
		}
	}

	return nil
}

func (k *KafkaConsumer) Stop(ctx context.Context) error {
	if k.consumer != nil {
		return k.consumer.Close()
	}

	return nil
}
