package clients

import (
	"context"

	"test-task-photo-booth/src/entities/dtos"
)

type PhotoStorage interface {
	Create(ctx context.Context, photo *dtos.PhotoDB) error
	FindOne(ctx context.Context, id string) (dtos.PhotoDB, error)
	Update(ctx context.Context, photo dtos.PhotoDB) error
	Delete(ctx context.Context, id string) error
}

type PhotoQueue interface {
	Publish(photo *dtos.Photo) error
}
