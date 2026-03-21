package kafkamodule

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader  *kafka.Reader
	logger  *zap.SugaredLogger
	handler Handler
}

type Handler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}

func New(brokers []string, topic, groupID string, handler Handler, logger *zap.SugaredLogger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})
	return &Consumer{
		reader:  reader,
		handler: handler,
		logger:  logger.With("component", "kafka-consumer", "topic", topic),
	}
}

func (c *Consumer) Run(ctx context.Context) {
    c.logger.Info("consumer started")
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                c.logger.Info("consumer stopped")
                return
            }
            c.logger.Errorw("fetch message failed", "err", err)
            continue
        }
        c.logger.Infow("message fetched", "topic", msg.Topic, "offset", msg.Offset) // <- добавь
        if err := c.handler.Handle(ctx, msg); err != nil {
            c.logger.Errorw("handle message failed", "err", err)
            continue
        }
        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            c.logger.Errorw("commit message failed", "err", err)
        }
    }
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
