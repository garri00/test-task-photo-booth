package rabbitmq

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"test-task-photo-booth/api/adapters/db/postgres"
	"test-task-photo-booth/api/usecases"
	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"
	"test-task-photo-booth/src/entities/dtos"
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

	log.Info().Msg("successfully connected to RabbitMQ")

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

const restartTimer = 20

func (c *RabbitMqClient) Listen(postgresClient *pgxpool.Pool, log *zerolog.Logger) error {
	//Add queues listeners
	for {
		err := c.photoQueue(postgresClient, log)
		if err != nil {
			log.Error().Err(err).Msgf("c.photoQueue failed; RESTART in %v sec", restartTimer)
			time.Sleep(time.Second * restartTimer)
		}
	}
}

func (c *RabbitMqClient) photoQueue(postgresClient *pgxpool.Pool, log *zerolog.Logger) error {
	defer func() {
		if r := recover(); r != nil {
			log.Info().Msgf("recovered from panic: %v; RESTART in %v sec", r, restartTimer)
			time.Sleep(time.Second * restartTimer)
		}
	}()

	photoCollection := postgres.NewPhotoStoragePG(postgresClient, log)
	photoUseCase := usecases.NewPhotoConsumeUseCase(photoCollection, log)

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

	var forever chan struct{}

	go func() {
		for m := range messages {
			if err := photoUseCase.Create(&dtos.Photo{Data: string(m.Body)}); err != nil {
				logger.Log.Error().Msgf("failed to create photo with body: %s", string(m.Body))
			}
		}
	}()

	<-forever

	return nil
}
