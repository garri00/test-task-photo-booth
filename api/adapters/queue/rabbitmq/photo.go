package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/clients"
	"test-task-photo-booth/src/entities"
	"test-task-photo-booth/src/entities/dtos"
)

type photoRmqQueue struct {
	ch     *amqp.Channel
	logger *zerolog.Logger
}

func NewPhotoRmqQueue(ch *amqp.Channel, logger *zerolog.Logger) clients.PhotoQueue {
	return &photoRmqQueue{
		ch:     ch,
		logger: logger,
	}
}

func (p photoRmqQueue) Publish(photo *dtos.Photo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()

	if err := p.ch.PublishWithContext(
		ctx,                  // Context for managing timeout
		"",                   // exchange (default exchange)
		entities.PhotosQueue, // routing key (queue name)
		false,                // mandatory (if true, the server will return an unroutable message)
		false,                // immediate (if true, the server will return an undeliverable message)
		amqp.Publishing{
			ContentType: "text/plain",       // Content type of the message
			Body:        []byte(photo.Data), // Message body as a byte array
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}
