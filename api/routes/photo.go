package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"test-task-photo-booth/api/adapters/db/postgres"
	"test-task-photo-booth/api/adapters/queue/rmq"
	"test-task-photo-booth/pkg/clients/rabbitmq"

	"test-task-photo-booth/api/handlers"
	"test-task-photo-booth/api/usecases"
)

func photo(postgresClient *pgxpool.Pool, rabbitClient *rabbitmq.RabbitMqClient, log *zerolog.Logger, r chi.Router) {
	photoCollection := postgres.NewPhotoStoragePG(postgresClient, log)
	photoQueue := rmq.NewPhotoProducer(rabbitClient.Conn, rabbitClient.PhotoQueue, log)
	photoUseCase := usecases.NewPhotoUseCase(photoCollection, photoQueue, log)
	photoHandler := handlers.NewPhotoHandler(photoUseCase, log)

	r.Post("/", photoHandler.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", photoHandler.GetByID)
		r.Delete("/", photoHandler.Delete)
	})

}
