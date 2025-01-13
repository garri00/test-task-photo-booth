package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"test-task-photo-booth/api"
	"test-task-photo-booth/api/adapters/db/postgres"
	"test-task-photo-booth/api/usecases"
	"test-task-photo-booth/pkg/clients/postgresql"
	"test-task-photo-booth/pkg/clients/rabbitmq"
	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/pkg/utils"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"

	rm "test-task-photo-booth/api/adapters/queue/rabbitmq"
)

func main() {
	//Init project and env configs
	configs, err := config.GetConfig()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to load .env")
	}

	//Setup main logger with level
	if err := logger.SetLogger(configs); err != nil {
		logger.Log.Fatal().Err(err).Msg("logger.SetMainLogger failed")
	}

	//Setup postgres DB connection
	ctx := context.Background()
	postgresClient, err := postgresql.NewClient(ctx, configs.PostgresConf, &logger.Log)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create postgres connection")
	}

	//Do all migrations
	if err := postgresql.MigrateUp(configs.PostgresConf, logger.Log); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to start postgres migration")
	}

	//Add queue
	rabbitmqClient, err := rabbitmq.NewRabbitMQ(configs.RabbitMQConf)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create rabbitmq connection")
	}

	ch, err := rabbitmqClient.Channel()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create create channel")
	}

	photoCollection := postgres.NewPhotoStoragePG(postgresClient, &logger.Log)
	photoQueue := rm.NewPhotoRmqQueue(ch, &logger.Log)
	photoUseCase := usecases.NewPhotoUseCase(photoCollection, photoQueue, &logger.Log)

	consumer, err := rabbitmq.NewConsumer(rabbitmqClient, photoUseCase)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create rabbitmq consumer")
	}

	if err := consumer.Listen(); err != nil {
		logger.Log.Fatal().Err(err).Msg("consumer.Listen() failed")
	}

	//Attach routes
	routes := api.NewRouter(postgresClient, ch, &logger.Log)

	logger.Log.Info().Msg("server started")

	//Set server attributes
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%s", configs.Host, configs.Port),
		Handler:           routes,
		ReadHeaderTimeout: entities.ServiceRequestTimeout * time.Second,
	}

	serviceVersion := "v1.0.0"

	//Start server
	logger.Log.Info().Msgf("service version=%s", serviceVersion)
	logger.Log.Info().Msgf("host=%s, port=%s, tls=%v, mode=%s", configs.Host, configs.Port, configs.TLS, configs.Mode)

	if configs.TLS {
		//Server with TLS
		serverTLSCert, err := utils.LoadCertificate()
		if err != nil {
			logger.Log.Fatal().Err(fmt.Errorf("utils.LoadCertificate() failed: %w", err)).Msg("Error loading certificate and key file")
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*serverTLSCert},
			MinVersion:   tls.VersionTLS12,
		}

		logger.Log.Debug().Msgf("https://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServeTLS("", ""); err != nil {
			logger.Log.Fatal().Err(err).Msg("server crashed")
		}
	} else {
		//Regular http server
		logger.Log.Debug().Msgf("http://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServe(); err != nil {
			logger.Log.Fatal().Err(err).Msg("server crashed")
		}
	}
}
