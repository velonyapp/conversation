
package transport

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

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
	groupID = "conversation"

	identityUserTopic = "identity.user"

	identityUserCreatedEvent = "identity.user.created"
	identityUserDeletedEvent = "identity.user.deleted"
)

var _ transport.Server = (*KafkaConsumer)(nil)

type KafkaConsumer struct {
	c   *conf.Transport
	svc *service.Consumer
	log *slog.Logger

	consumer *kafka.Consumer
}

func NewKafkaConsumer(
	c *conf.Transport,
	svc *service.Consumer,
	log *slog.Logger,
) *KafkaConsumer {
	return &KafkaConsumer{
		c:   c,
		svc: svc,
		log: log,
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
		k.log.ErrorContext(ctx,
			"failed to create Kafka consumer",
			slog.Any("error", err),
		)
		return err
	}

	k.consumer = consumer

	topics := []string{
		identityUserTopic,
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		k.log.ErrorContext(ctx,
			"failed to subscribe to Kafka topics",
			slog.Any("topics", topics),
			slog.Any("error", err),
		)
		return err
	}

	k.log.InfoContext(ctx,
		"Kafka consumer subscribed to topics",
		slog.String("group_id", groupID),
		slog.Any("topics", topics),
	)

	for !consumer.IsClosed() {
		event := consumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			retryDelay := k.c.Kafka.RetryInitialInterval.AsDuration()

			for attempt := 1; ; attempt++ {
				if err := k.handleMessage(ctx, e); err != nil && !errors.Is(err, ErrInvalidMessage) {
					k.log.WarnContext(ctx,
						"failed to handle Kafka message, retrying",
						slog.Int("attempt", attempt),
						slog.Duration("retry_delay", retryDelay),
						slog.Any("topic_partition", e.TopicPartition),
						slog.Any("error", err),
					)

					time.Sleep(retryDelay)

					retryDelay = min(
						time.Duration(float64(retryDelay)*k.c.Kafka.RetryMultiplier),
						k.c.Kafka.RetryMaxInterval.AsDuration(),
					)

					continue
				}

				break
			}

			consumer.StoreMessage(e)

		case kafka.Error:
			if e.IsFatal() {
				k.log.ErrorContext(ctx,
					"fatal Kafka consumer error",
					slog.Any("error", e),
				)
				return errors.New(e.Error())
			}

			k.log.WarnContext(ctx,
				"Kafka consumer error",
				slog.Any("error", e),
			)
		}
	}

	return nil
}

func (k *KafkaConsumer) Stop(ctx context.Context) error {
	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			k.log.ErrorContext(ctx,
				"failed to close Kafka consumer",
				slog.Any("error", err),
			)
			return err
		}

		k.log.InfoContext(ctx,
			"Kafka consumer stopped",
			slog.String("group_id", groupID),
		)
	}

	return nil
}
