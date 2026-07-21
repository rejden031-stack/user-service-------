package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type UserCreatedEvent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddr string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

func (p *Producer) PublishUserCreated(ctx context.Context, event UserCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", event.ID)),
			Value: data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
