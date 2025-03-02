package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"test-task-photo-booth/api"
	"test-task-photo-booth/pkg/clients/postgresql"
	"test-task-photo-booth/pkg/clients/rabbitmq"
	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/pkg/utils"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"
)

const loggerName = "producer"

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

	//Setup postgres DB connection
	ctx := context.Background()
	postgresClient, err := postgresql.NewClient(ctx, configs.PostgresConf, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create postgres connection")
	}

	//Do all migrations
	if err := postgresql.MigrateUp(configs.PostgresConf, log); err != nil {
		log.Fatal().Err(err).Msg("failed to start postgres migration")
	}

	//Add queues
	rabbitmqConn, err := rabbitmq.NewRabbitMqConnection(configs.RabbitMQConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rabbitmq connection")
	}

	rabbitmqClient, err := rabbitmq.NewRabbitMqClient(rabbitmqConn, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rabbitmq consumer")
	}

	//Attach routes
	routes := api.NewRouter(postgresClient, rabbitmqClient, log)

	log.Info().Msg("server started")

	//Set server attributes
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%s", configs.Host, configs.Port),
		Handler:           routes,
		ReadHeaderTimeout: entities.ServiceRequestTimeout * time.Second,
	}

	serviceVersion := "v1.0.0"

	//Start server
	log.Info().Msgf("service version=%s", serviceVersion)
	log.Info().Msgf("host=%s, port=%s, tls=%v, mode=%s", configs.Host, configs.Port, configs.TLS, configs.Mode)

	if configs.TLS {
		//Server with TLS
		serverTLSCert, err := utils.LoadCertificate()
		if err != nil {
			log.Fatal().Err(fmt.Errorf("utils.LoadCertificate() failed: %w", err)).Msg("Error loading certificate and key file")
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*serverTLSCert},
			MinVersion:   tls.VersionTLS12,
		}

		log.Debug().Msgf("https://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal().Err(err).Msg("server crashed")
		}
	} else {
		//Regular http server
		log.Debug().Msgf("http://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("server crashed")
		}
	}
}
