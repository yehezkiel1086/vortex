package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type EventPublisher struct {
	mq *RabbitMQ
}

func NewEventPublisher(mq *RabbitMQ) port.EventPublisher {
	return &EventPublisher{
		mq: mq,
	}
}

func (p *EventPublisher) PublishEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return domain.ErrInvalidInput
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		MessageId:    event.ID.String(),
		Body:         body,
	}

	err = p.mq.Channel.PublishWithContext(
		ctx,
		ExchangeEvents,
		RoutingKeyIngested,
		false, // mandatory
		false, // immediate
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to rabbitmq: %w", err)
	}

	return nil
}

// Ensure EventPublisher implements port.EventPublisher
var _ port.EventPublisher = (*EventPublisher)(nil)
