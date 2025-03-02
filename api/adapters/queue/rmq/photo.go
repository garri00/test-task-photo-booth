package rmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/clients"
	"test-task-photo-booth/pkg/clients/rabbitmq"
	"test-task-photo-booth/src/entities/dtos"
)

type photoProducer struct {
	client     *amqp.Connection
	photoQueue rabbitmq.PhotoQueue
	logger     *zerolog.Logger
}

func NewPhotoProducer(client *amqp.Connection, photoQueue rabbitmq.PhotoQueue, logger *zerolog.Logger) clients.PhotoQueue {
	return &photoProducer{
		client:     client,
		photoQueue: photoQueue,
		logger:     logger,
	}
}

func (p photoProducer) Publish(photo *dtos.Photo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()

	ch, err := p.client.Channel()
	if err != nil {
		return fmt.Errorf("could not open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.PublishWithContext(
		ctx,               // Context for managing timeout
		"",                // exchange (default exchange)
		p.photoQueue.Name, // routing key (queue name)
		false,             // mandatory (if true, the server will return an unroutable message)
		false,             // immediate (if true, the server will return an undeliverable message)
		amqp.Publishing{
			ContentType: "text/plain",       // Content type of the message
			Body:        []byte(photo.Data), // Message body as a byte array
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}
