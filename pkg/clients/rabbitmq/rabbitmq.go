package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"test-task-photo-booth/api/usecases"
	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"
	"test-task-photo-booth/src/entities/dtos"
)

func NewRabbitMQ(configs config.RabbitMQConf) (*amqp.Connection, error) {
	connectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		configs.Username,
		configs.Password,
		configs.Host,
		configs.Port)

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}

type Consumer struct {
	conn         *amqp.Connection
	photoUseCase usecases.PhotoUseCase
	log          *zerolog.Logger
}

func NewConsumer(conn *amqp.Connection, photoUseCase usecases.PhotoUseCase) (Consumer, error) {
	consumer := Consumer{conn: conn, photoUseCase: photoUseCase}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, fmt.Errorf("failed to setup consumer: %w", err)
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	_, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	return nil
}

func (c *Consumer) Listen() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	q, err := declareQueue(ch, entities.PhotosQueue)
	if err != nil {
		return fmt.Errorf("failed to declare photo queue: %w", err)
	}

	messages, err := ch.Consume(
		q.Name, // Queue name
		"",     // consumer tag (empty string means a unique tag will be generated)
		true,   // auto-ack (automatically acknowledge messages)
		false,  // exclusive (only this consumer can access the queue)
		false,  // no-local (if true, the server will not deliver messages to the connection that published them)
		false,  // no-wait (do not wait for a server response)
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// Create a channel that runs in the background forever
	go func(photoUseCase usecases.PhotoUseCase) {
		for {
			for m := range messages {
				if err := photoUseCase.Create(&dtos.Photo{Data: string(m.Body)}); err != nil {
					logger.Log.Error().Msgf("failed to create photo with body: %s", string(m.Body))
				}
			}
		}

	}(c.photoUseCase)

	return nil
}

func declareQueue(ch *amqp.Channel, name string) (*amqp.Queue, error) {
	//Declares all queues
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &q, nil
}
