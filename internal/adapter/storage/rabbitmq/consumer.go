package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type EventHandler func(ctx context.Context, event *domain.Event) error

type Consumer struct {
	mq *RabbitMQ
}

func NewConsumer(mq *RabbitMQ) *Consumer {
	return &Consumer{mq: mq}
}

func (c *Consumer) StartConsuming(ctx context.Context, handler EventHandler) error {
	// Set prefetch count for fair dispatch
	if err := c.mq.Channel.Qos(10, 0, false); err != nil {
		return fmt.Errorf("failed to set rabbitmq qos: %w", err)
	}

	msgs, err := c.mq.Channel.Consume(
		QueueEventsIngested,
		"vortex-worker", // consumer tag
		false,           // auto-ack (false so we ack after successful processing)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming from queue: %w", err)
	}

	log.Printf("[Worker] Started consuming from queue '%s'...", QueueEventsIngested)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Worker] Stopping consumer loop...")
			return nil
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq delivery channel closed")
			}

			var event domain.Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("[Worker] Failed to unmarshal message: %v. Rejecting message.", err)
				_ = d.Nack(false, false) // discard invalid message
				continue
			}

			if err := handler(ctx, &event); err != nil {
				log.Printf("[Worker] Handler error processing event %s: %v. Re-queuing.", event.ID, err)
				_ = d.Nack(false, true) // requeue for retry
			} else {
				_ = d.Ack(false)
			}
		}
	}
}
