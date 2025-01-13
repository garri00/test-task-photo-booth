package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"test-task-photo-booth/api/adapters/db/postgres"
	"test-task-photo-booth/api/adapters/queue/rabbitmq"
	"test-task-photo-booth/api/handlers"
	"test-task-photo-booth/api/usecases"
)

func photo(postgresClient *pgxpool.Pool, ch *amqp.Channel, log *zerolog.Logger, r chi.Router) {
	photoCollection := postgres.NewPhotoStoragePG(postgresClient, log)
	photoQueue := rabbitmq.NewPhotoRmqQueue(ch, log)
	photoUseCase := usecases.NewPhotoUseCase(photoCollection, photoQueue, log)
	photoHandler := handlers.NewPhotoHandler(photoUseCase, log)

	r.Post("/", photoHandler.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", photoHandler.GetByID)
		r.Delete("/", photoHandler.Delete)
	})

}
