package main

import (
	"context"

	"test-task-photo-booth/pkg/clients/postgresql"
	"test-task-photo-booth/pkg/clients/rabbitmq"
	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/config"
)

const loggerName = "consumer"

func main() {
	//Init project and env configs
	configs, err := config.GetConfig()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to load .env")
	}

	//Setup main logger with level
	log, err := logger.SetServiceLogger(loggerName, configs)
	if err != nil {
		log.Fatal().Err(err).Msg("logger.SetMainLogger failed")
	}

	ctx := context.Background()
	postgresClient, err := postgresql.NewClient(ctx, configs.PostgresConf, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create postgres connection")
	}

	//Add queue
	rabbitConn, err := rabbitmq.NewRabbitMqConnection(configs.RabbitMQConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rabbitmq connection")
	}

	consumer, err := rabbitmq.NewRabbitMqClient(rabbitConn, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rabbitmq consumer")
	}

	if err := consumer.Listen(postgresClient, log); err != nil {
		log.Fatal().Err(err).Msg("consumer.Listen() failed")
	}

	serviceVersion := "v1.0.0"

	//Start server
	log.Info().Msgf("consumer version=%s", serviceVersion)
	log.Info().Msg("consumer started successfully")
}
