package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/clients/rabbitmq"
)

func API(postgresClient *pgxpool.Pool, rabbitClient *rabbitmq.RabbitMqClient, log *zerolog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Route("/photo", func(r chi.Router) {
		photo(postgresClient, rabbitClient, log, r)
	})

	return r
}
