package kafkamodule

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Config struct {
	Brokers []string
}

type Publisher struct {
	writer *kafka.Writer
	logger *zap.SugaredLogger
}

func NewPublisher(cfg Config, logger *zap.SugaredLogger) (*Publisher, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Publisher{
		writer: writer,
		logger: logger.With("component", "kafka-publisher"),
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	msg := kafka.Message{
		Topic: eventType, // каждый eventType — отдельный топик
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(eventType)},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka write message: %w", err)
	}

	p.logger.Debugw("event published", "event_type", eventType)
	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
