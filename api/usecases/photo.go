package usecases

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/clients"
	"test-task-photo-booth/pkg/utils"
	"test-task-photo-booth/src/entities/dtos"
)

type PhotoUseCase struct {
	db    clients.PhotoStorage
	queue clients.PhotoQueue
	log   *zerolog.Logger
}

func NewPhotoUseCase(storage clients.PhotoStorage, queue clients.PhotoQueue, l *zerolog.Logger) PhotoUseCase {
	return PhotoUseCase{
		db:    storage,
		queue: queue,
		log:   l,
	}
}

const (
	photoResize75 = 75
	photoResize50 = 50
	photoResize25 = 25
)

func (p PhotoUseCase) AddInQueue(photo *dtos.Photo) error {
	if err := p.queue.Publish(photo); err != nil {
		return fmt.Errorf("error adding photo to queue: %v", err)
	}

	return nil
}

func (p PhotoUseCase) Create(photo *dtos.Photo) error {
	ctx := context.Background()

	decodeData, err := base64.StdEncoding.DecodeString(photo.Data)
	if err != nil {
		return fmt.Errorf("could not decode data: %w", err)
	}

	extension := utils.GetB64MimeType(decodeData)

	data75, err := utils.ResizeImageB64(photo.Data, extension, photoResize75)
	if err != nil {
		return fmt.Errorf("utils.ResizeImageB64() failed: %w", err)
	}

	data50, err := utils.ResizeImageB64(photo.Data, extension, photoResize50)
	if err != nil {
		return fmt.Errorf("utils.ResizeImageB64() failed: %w", err)
	}

	data25, err := utils.ResizeImageB64(photo.Data, extension, photoResize25)
	if err != nil {
		return fmt.Errorf("utils.ResizeImageB64() failed: %w", err)
	}

	photoDB := &dtos.PhotoDB{
		DataOrigin: photo.Data,
		Data75:     data75,
		Data50:     data50,
		Data25:     data25,
		IsDeleted:  false,
	}

	if err := p.db.Create(ctx, photoDB); err != nil {
		return fmt.Errorf("db.Create(): %w", err)
	}

	photo.ID = photoDB.ID
	photo.IsDeleted = photoDB.IsDeleted

	return nil
}

func (p PhotoUseCase) GetByID(id, quality string) (dtos.Photo, error) {
	ctx := context.Background()
	photoDB, err := p.db.FindOne(ctx, id)
	if err != nil {
		return dtos.Photo{}, fmt.Errorf("db.FindOne(): %w", err)
	}

	photo := p.getPhotoWithQuality(photoDB, quality)

	return photo, nil
}

func (p PhotoUseCase) getPhotoWithQuality(photoDB dtos.PhotoDB, quality string) dtos.Photo {
	photo := dtos.Photo{ID: photoDB.ID, IsDeleted: photoDB.IsDeleted}

	switch quality {
	case "100":
		photo.Data = photoDB.DataOrigin
	case "75":
		photo.Data = photoDB.Data75
	case "50":
		photo.Data = photoDB.Data50
	case "25":
		photo.Data = photoDB.Data25
	default:
		p.log.Warn().Msgf("PhotoUseCase getPhotoWithQuality failed: photo quality:'%s' not defined, original photo provided", quality)
		photo.Data = photoDB.DataOrigin
	}

	return photo
}

func (p PhotoUseCase) Delete(id string) error {
	ctx := context.Background()
	if err := p.db.Delete(ctx, id); err != nil {
		return fmt.Errorf("db.Delete(): %w", err)
	}

	return nil
}
