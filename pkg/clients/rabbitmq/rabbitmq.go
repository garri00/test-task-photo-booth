package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"
)

func NewRabbitMqConnection(configs config.RabbitMQConf) (*amqp.Connection, error) {
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

type PhotoQueue *amqp.Queue

type RabbitMqClient struct {
	Conn       *amqp.Connection
	PhotoQueue PhotoQueue
	log        *zerolog.Logger
}

func NewRabbitMqClient(conn *amqp.Connection, log *zerolog.Logger) (*RabbitMqClient, error) {
	rabbitMq := RabbitMqClient{
		Conn: conn,
		log:  log,
	}

	if err := rabbitMq.setupChannels(); err != nil {
		return nil, fmt.Errorf("failed to setup consumer: %w", err)
	}

	return &rabbitMq, nil
}

func (c *RabbitMqClient) setupChannels() error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	photoQueue, err := declareQueue(ch, entities.PhotosQueue)
	if err != nil {
		return fmt.Errorf("failed to declare photo queue: %w", err)
	}

	c.PhotoQueue = photoQueue

	return nil
}

func (c *RabbitMqClient) Listen() error {
	//photoCollection := postgres.NewPhotoStoragePG(postgresClient, log)
	//photoUseCase := usecases.NewPhotoUseCase(photoCollection, c.PhotoQueue, log)
	//photoHandler := handlers.NewPhotoHandler(photoUseCase, log)

	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	messages, err := ch.Consume(
		c.PhotoQueue.Name, // Queue name
		"",                // consumer tag (empty string means a unique tag will be generated)
		true,              // auto-ack (automatically acknowledge messages)
		false,             // exclusive (only this consumer can access the queue)
		false,             // no-local (if true, the server will not deliver messages to the connection that published them)
		false,             // no-wait (do not wait for a server response)
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	//Select queus defer panic

	// Create a channel that runs in the background forever
	go func() {
		for {
			for m := range messages {
				//if err := photoUseCase.Create(&dtos.Photo{Data: string(m.Body)}); err != nil {
				//	logger.Log.Error().Msgf("failed to create photo with body: %s", string(m.Body))
				//}
				logger.Log.Error().Msgf("failed to create photo with body: %s", string(m.Body))
			}
		}
	}()

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
