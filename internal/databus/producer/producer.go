package producer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	producer *kafka.Writer
}

func New() *Producer {
	producer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "user.register",
		BatchSize:    0,
		RequiredAcks: kafka.RequireAll,
		WriteTimeout: 10 * time.Second,
	}
	return &Producer{
		producer: producer,
	}
}

func (p *Producer) SendMessage(ctx context.Context, msg []byte, key []byte) error {
	err := p.producer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: msg,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}
	return nil
}
