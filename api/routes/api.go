package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func API(postgresClient *pgxpool.Pool, ch *amqp.Channel, log *zerolog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Route("/photo", func(r chi.Router) {
		photo(postgresClient, ch, log, r)
	})

	return r
}
